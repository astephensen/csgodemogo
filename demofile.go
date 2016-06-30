package csgodemogo

import (
	"fmt"
	"os"

	"github.com/astephensen/csgodemogo/cstrikeproto"
	"github.com/golang/protobuf/proto"
)

type demoFile struct {
	Header        demoHeader
	Tick          int32
	Frame         int32
	Finished      bool
	stream        *demoStream
	gameEventList cstrikeproto.CSVCMsg_GameEventList
	// Emitter Functions - these should probably be channels
	GameEventListEmitter func(gameEventList *cstrikeproto.CSVCMsg_GameEventList)
	GameEventEmitter     func(gameEvent interface{})
}

func Open(path string) *demoFile {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	demo := demoFile{}
	demo.stream = DemoStream(file)
	demo.Header = DemoHeader(demo.stream)
	demo.Tick = 0
	demo.Frame = 0
	demo.Finished = false

	return &demo
}

func (demo *demoFile) GetFrame() {
	commandHeader := CommandHeader(demo.stream)
	switch commandHeader.Command {
	case CommandSignOn, CommandPacket:
		// First 160 bytes consist of the command info and sequence number in / out
		// Let's just ignore this for now.
		demo.stream.Skip(160)
		dataLength := int(demo.stream.GetInt())
		demo.ParseProtobufPacket(dataLength)

	case CommandSyncTick:
		// Sync ticks will be ignored.

	case CommandConsoleCommand:

	case CommandUserCommand:

	case CommandDataTables:
		// Skip over data tables.
		dataLength := demo.stream.GetInt()
		demo.stream.Skip(int64(dataLength))

	case CommandStop:
		demo.Finished = true

	case CommandCustomData:

	case CommandStringTables:
		// Skip over string tables.
		dataLength := demo.stream.GetInt()
		demo.stream.Skip(int64(dataLength))
	}
}

func (demo *demoFile) ParseProtobufPacket(length int) {
	// The packet can contain multiple commands so loop while the length is there.
	currentPosition := demo.stream.pos
	for demo.stream.pos < currentPosition+length {
		command := demo.stream.GetVarInt()
		size := demo.stream.GetVarInt()

		// Read the body into a buffer.
		buffer := make([]byte, size)
		_, err := demo.stream.Read(buffer)
		if err != nil {
			panic(err)
		}

		if command <= 7 {
			// NET messages have commands between 0 and 7
			//netMessage := cstrikeproto.NET_Messages(command)
			//fmt.Println("Got net message", netMessage)
		} else if command <= 31 {
			// SVC messages have commands between 8 and 31
			svcMessage := cstrikeproto.SVC_Messages(command)
			//fmt.Println("Got svc message", svcMessage)

			switch svcMessage {

			// Parse the game event table.
			case cstrikeproto.SVC_Messages_svc_GameEventList:
				err := proto.Unmarshal(buffer, &demo.gameEventList)
				if err != nil {
					panic(err)
				}
				if demo.GameEventListEmitter != nil {
					demo.GameEventListEmitter(&demo.gameEventList)
				}

			// Parse game events into a user friendly version.
			case cstrikeproto.SVC_Messages_svc_GameEvent:
				if demo.GameEventEmitter != nil {
					gameEvent := ParseGameEvent(&demo.gameEventList, buffer)
					if gameEvent != nil {
						demo.GameEventEmitter(gameEvent)
					}
				}
			}

		} else {
			// All other commands are currently unknown.
		}
	}
}
