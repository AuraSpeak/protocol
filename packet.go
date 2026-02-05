package protocol

import (
	"bytes"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
)

const (
	Magic   = 0xAE
	Version = 0x01
)

// HeaderSize is the size of the header in bytes
// It is the size of the packet type
const HeaderSize = 14

// Header is the header of the packet
// It contains the packet type
type Header struct {
	Magic      uint8
	Version    uint8
	PacketType PacketType
	Flags      uint8
	ID         uint32
	Length     uint32
	FragmentID uint16
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
		return nil, ErrDataTooShort
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
		return Header{}, ErrDataTooShort
	}
	magic := data[0]
	version := data[1]
	packetType := PacketType(data[2])
	flags := data[3]
	id := binary.BigEndian.Uint32(data[4:8])
	length := binary.BigEndian.Uint32(data[8:12])
	fragmentID := binary.BigEndian.Uint16(data[12:14])

	if !IsValidPacketType(packetType) {
		return Header{}, ErrInvalidPacketType
	}
	if magic != Magic {
		return Header{}, ErrInvalidMagic
	}
	if version != Version {
		return Header{}, ErrInvalidVersion
	}

	header := Header{
		Magic:      magic,
		Version:    version,
		PacketType: packetType,
		Flags:      flags,
		ID:         id,
		Length:     length,
		FragmentID: fragmentID,
	}
	log.WithField("caller", "protocol").Infof("Decoded header: %s", PacketTypeMapType[packetType])
	return header, nil
}

// EncodeHeader encodes the header into a byte slice
// Example:
//
//	header := Header{PacketType: PacketTypeDebugHello}
//	encoded := EncodeHeader(header)
//	fmt.Println(encoded)
func EncodeHeader(header Header) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, HeaderSize))
	buf.WriteByte(byte(header.Magic))
	buf.WriteByte(byte(header.Version))
	buf.WriteByte(byte(header.PacketType))
	buf.WriteByte(byte(header.Flags))
	wirteUint32(buf, header.ID)
	wirteUint32(buf, header.Length)
	wirteUint16(buf, header.FragmentID)
	return buf.Bytes()
}

func wirteUint16(buf *bytes.Buffer, value uint16) {
	binary.Write(buf, binary.BigEndian, value)
}

func wirteUint32(buf *bytes.Buffer, value uint32) {
	binary.Write(buf, binary.BigEndian, value)
}

func readUint16(buf *bytes.Buffer) uint16 {
	var value uint16
	binary.Read(buf, binary.BigEndian, &value)
	return value
}

func readUint32(buf *bytes.Buffer) uint32 {
	var value uint32
	binary.Read(buf, binary.BigEndian, &value)
	return value
}
