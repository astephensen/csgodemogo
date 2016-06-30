package csgodemogo

const (
	CommandSignOn         = 1
	CommandPacket         = 2
	CommandSyncTick       = 3
	CommandConsoleCommand = 4
	CommandUserCommand    = 5
	CommandDataTables     = 6
	CommandStop           = 7
	CommandCustomData     = 8
	CommandStringTables   = 9
	CommandLastCommand    = 9
)

type demoCommandHeader struct {
	Command    byte
	Tick       int32
	PlayerSlot byte
}

func CommandHeader(stream *demoStream) demoCommandHeader {
	commandHeader := demoCommandHeader{
		Command:    stream.GetByte(),
		Tick:       stream.GetInt(),
		PlayerSlot: stream.GetByte(),
	}
	return commandHeader
}
