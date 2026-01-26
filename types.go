package protocol

//go:generate go run ./cmd/generate
//go:generate go run golang.org/x/tools/cmd/stringer -type=PacketType -trimprefix=PacketType

type PacketType uint8

type PacketTypeMapping struct {
	PacketType PacketType
	String     string
}

// IsValidPacketType checks if the packet type is valid
// It returns true if the packet type is valid, false otherwise
func IsValidPacketType(packetType PacketType) bool {
	return packetType != PacketTypeNone
}
