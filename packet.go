package protocol

import (
	"bytes"
	"errors"

	log "github.com/sirupsen/logrus"
)

// HeaderSize is the size of the header in bytes
// It is the size of the packet type
const HeaderSize = 1

// Header is the header of the packet
// It contains the packet type
type Header struct {
	PacketType PacketType
}

// Packet is the packet of the protocol
// It contains the header and the payload
type Packet struct {
	PacketHeader Header
	Payload      []byte
}

// Encode encodes the packet into a byte slice
// Example:
//
//	packet := &Packet{
//		PacketHeader: Header{PacketType: PacketTypeDebugHello},
//		Payload: []byte("Hello, Server!"),
//	}
//	encoded := packet.Encode()
//	fmt.Println(encoded)
func (p *Packet) Encode() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, len(p.Payload)+HeaderSize))
	buf.Write(EncodeHeader(p.PacketHeader))
	buf.Write(p.Payload)
	return buf.Bytes()
}

// Decode decodes the packet from a byte slice
// Example:
//
//	packet, err := Decode(Header{PacketType: PacketTypeDebugHello}, []byte("Hello, Server!"))
//	if err != nil {
//		fmt.Println("Error decoding packet:", err)
//	}
//	fmt.Println(packet.Payload)
func Decode(data []byte) (*Packet, error) {
	if len(data) < HeaderSize {
		return nil, errors.New("data too short")
	}
	packetHeader, err := DecodeHeader(data[:HeaderSize])
	if err != nil {
		return nil, err
	}
	payload := data[HeaderSize:]
	return &Packet{
		PacketHeader: packetHeader,
		Payload:      payload,
	}, nil
}

// DecodeHeader decodes the header from a byte slice
// Example:
//
//	header, err := DecodeHeader([]byte{0x01})
//	if err != nil {
//		fmt.Println("Error decoding header:", err)
//	}
//	fmt.Println(header.PacketType)
func DecodeHeader(data []byte) (Header, error) {
	if len(data) < HeaderSize {
		return Header{}, errors.New("data too short")
	}
	packetType := PacketType(data[0])
	if !IsValidPacketType(packetType) {
		return Header{}, errors.New("invalid packet type")
	}
	log.WithField("caller", "protocol").Infof("Decoded header: %s", PacketTypeMapType[packetType])
	return Header{PacketType: packetType}, nil
}

// EncodeHeader encodes the header into a byte slice
// Example:
//
//	header := Header{PacketType: PacketTypeDebugHello}
//	encoded := EncodeHeader(header)
//	fmt.Println(encoded)
func EncodeHeader(header Header) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, HeaderSize))
	buf.WriteByte(byte(header.PacketType))
	return buf.Bytes()
}
