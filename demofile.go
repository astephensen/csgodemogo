package csgodemogo

import (
	"fmt"
	"os"

	"github.com/astephensen/csgodemogo/cstrikeproto"
	"github.com/golang/protobuf/proto"
)

type demoFile struct {
	Header demoHeader
	tick   int32
	frame  int32
	stream *demoStream
}

func Open(path string) *demoFile {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	demo := demoFile{}
	demo.stream = DemoStream(file)
	demo.Header = DemoHeader(demo.stream)
	demo.tick = 0
	demo.frame = 0

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
		fmt.Println("Sync Tick")

	case CommandConsoleCommand:
		fmt.Println("Console Command")

	case CommandUserCommand:
		fmt.Println("User Command")

	case CommandDataTables:
		fmt.Println("Data Tables")
		dataLength := demo.stream.GetInt()
		demo.stream.Skip(int64(dataLength))

	case CommandStop:
		fmt.Println("Stop!")

	case CommandCustomData:
		fmt.Println("Custom Data")

	case CommandStringTables:
		fmt.Println("String Tables")
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

			case cstrikeproto.SVC_Messages_svc_ServerInfo:
				message := cstrikeproto.CSVCMsg_ServerInfo{}
				err := proto.Unmarshal(buffer, &message)
				if err != nil {
					panic(err)
				}
				fmt.Println(message.String())

			case cstrikeproto.SVC_Messages_svc_GameEventList:
				message := cstrikeproto.CSVCMsg_GameEventList{}
				err := proto.Unmarshal(buffer, &message)
				if err != nil {
					panic(err)
				}
				fmt.Println(message.String())

			case cstrikeproto.SVC_Messages_svc_GameEvent:
				message := cstrikeproto.CSVCMsg_GameEvent{}
				err := proto.Unmarshal(buffer, &message)
				if err != nil {
					panic(err)
				}
				fmt.Println(message.String())

				/*
					case cstrikeproto.SVC_Messages_svc_PacketEntities:
						message := cstrikeproto.CSVCMsg_PacketEntities{}
						err := proto.Unmarshal(buffer, &message)
						if err != nil {
							panic(err)
						}
						fmt.Println(message.String())
				*/

			}

		} else {
			// All other commands are currently unknown.
		}
	}
}
