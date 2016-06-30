package csgodemogo

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type demoHeader struct {
	DemoFilestamp   [8]byte
	DemoProtocol    int32
	NetworkProtocol int32
	ServerName      [MAX_OSPATH]byte
	ClientName      [MAX_OSPATH]byte
	MapName         [MAX_OSPATH]byte
	GameDirectory   [MAX_OSPATH]byte
	PlaybackTime    float32
	PlaybackTicks   int32
	PlaybackFrames  int32
	SignonLength    int32
}

// Return a DemoHeader from a demo stream.
func DemoHeader(stream *demoStream) demoHeader {
	header := demoHeader{}
	err := binary.Read(stream, binary.LittleEndian, &header)
	if err != nil {
		panic(err)
	}
	return header
}

// PrintInfo will print the header information.
func (header *demoHeader) PrintInfo() {
	fmt.Printf("Filestamp: %s\n", strings.TrimRight(string(header.DemoFilestamp[:]), "\x00"))
	fmt.Printf("Protocol: %d\n", header.DemoProtocol)
	fmt.Printf("Network Protocol: %d\n", header.NetworkProtocol)
	fmt.Printf("Server Name: %s\n", strings.TrimRight(string(header.ServerName[:]), "\x00"))
	fmt.Printf("Client name: %s\n", strings.TrimRight(string(header.ClientName[:]), "\x00"))
	fmt.Printf("Map: %s\n", strings.TrimRight(string(header.MapName[:]), "\x00"))
	fmt.Printf("Game Directory: %s\n", strings.TrimRight(string(header.GameDirectory[:]), "\x00"))
	fmt.Printf("Playback time: %f seconds\n", header.PlaybackTime)
	fmt.Printf("Ticks: %d\n", header.PlaybackTicks)
	fmt.Printf("Frames: %d\n", header.PlaybackFrames)
	fmt.Printf("Signon Length: %d\n", header.SignonLength)
}
