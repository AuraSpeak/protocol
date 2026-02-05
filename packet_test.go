package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

// buildHeaderBytes builds a 14-byte header slice matching DecodeHeader read order:
// bytes 0-3: magic, version, packetType, flags; 4-7: id (BigEndian); 8-11: length (BigEndian); 12-13: fragmentID (BigEndian).
func buildHeaderBytes(magic, version, packetType, flags uint8, id uint32, length uint32, fragmentID uint16) []byte {
	b := make([]byte, 14)
	b[0] = magic
	b[1] = version
	b[2] = packetType
	b[3] = flags
	binary.BigEndian.PutUint32(b[4:8], id)
	binary.BigEndian.PutUint32(b[8:12], length)
	binary.BigEndian.PutUint16(b[12:14], fragmentID)
	return b
}

func TestDecodeHeader_DataTooShort(t *testing.T) {
	for _, data := range [][]byte{
		nil,
		{},
		{0x01},
		buildHeaderBytes(Magic, Version, byte(PacketTypeDebugHello), 0, 0, 0, 0)[:13],
	} {
		_, err := DecodeHeader(data)
		if err == nil {
			t.Errorf("DecodeHeader(%v): expected error, got nil", data)
			continue
		}
		if !errors.Is(err, ErrDataTooShort) {
			t.Errorf("DecodeHeader(%v): expected ErrDataTooShort, got %v", data, err)
		}
	}
}

func TestDecodeHeader_InvalidMagic(t *testing.T) {
	data := buildHeaderBytes(Magic+1, Version, byte(PacketTypeDebugHello), 0, 0, 0, 0)
	_, err := DecodeHeader(data)
	if err == nil {
		t.Fatal("DecodeHeader(invalid magic): expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidMagic) {
		t.Errorf("DecodeHeader(invalid magic): expected ErrInvalidMagic, got %v", err)
	}
}

func TestDecodeHeader_InvalidVersion(t *testing.T) {
	data := buildHeaderBytes(Magic, Version+1, byte(PacketTypeDebugHello), 0, 0, 0, 0)
	_, err := DecodeHeader(data)
	if err == nil {
		t.Fatal("DecodeHeader(invalid version): expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidVersion) {
		t.Errorf("DecodeHeader(invalid version): expected ErrInvalidVersion, got %v", err)
	}
}

func TestDecodeHeader_InvalidPacketType(t *testing.T) {
	data := buildHeaderBytes(Magic, Version, byte(PacketTypeNone), 0, 0, 0, 0)
	_, err := DecodeHeader(data)
	if err == nil {
		t.Fatal("DecodeHeader(invalid packet type): expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidPacketType) {
		t.Errorf("DecodeHeader(invalid packet type): expected ErrInvalidPacketType, got %v", err)
	}
}

func TestDecodeHeader_Valid(t *testing.T) {
	magic := uint8(Magic)
	version := uint8(Version)
	packetType := PacketTypeDebugHello
	flags := uint8(0x42)
	id := uint32(1)
	length := uint32(100)
	fragmentID := uint16(2)
	data := buildHeaderBytes(magic, version, byte(packetType), flags, id, length, fragmentID)
	header, err := DecodeHeader(data)
	if err != nil {
		t.Fatalf("DecodeHeader(valid): unexpected error: %v", err)
	}
	if header.Magic != magic {
		t.Errorf("header.Magic = %v, want %v", header.Magic, magic)
	}
	if header.Version != version {
		t.Errorf("header.Version = %v, want %v", header.Version, version)
	}
	if header.PacketType != packetType {
		t.Errorf("header.PacketType = %v, want %v", header.PacketType, packetType)
	}
	if header.Flags != flags {
		t.Errorf("header.Flags = %v, want %v", header.Flags, flags)
	}
	if header.ID != id {
		t.Errorf("header.ID = %v, want %v", header.ID, id)
	}
	if header.Length != length {
		t.Errorf("header.Length = %v, want %v", header.Length, length)
	}
	if header.FragmentID != fragmentID {
		t.Errorf("header.FragmentID = %v, want %v", header.FragmentID, fragmentID)
	}
}

func TestDecode_DataTooShort(t *testing.T) {
	for _, data := range [][]byte{
		nil,
		{},
		make([]byte, 11),
	} {
		_, err := Decode(data)
		if err == nil {
			t.Errorf("Decode(len=%d): expected error, got nil", len(data))
			continue
		}
		if !errors.Is(err, ErrDataTooShort) {
			t.Errorf("Decode(len=%d): expected ErrDataTooShort, got %v", len(data), err)
		}
	}
}

func TestDecode_InvalidMagic(t *testing.T) {
	data := append(buildHeaderBytes(Magic+1, Version, byte(PacketTypeDebugHello), 0, 0, 0, 0), []byte("payload")...)
	_, err := Decode(data)
	if err == nil {
		t.Fatal("Decode(invalid magic): expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidMagic) {
		t.Errorf("Decode(invalid magic): expected ErrInvalidMagic, got %v", err)
	}
}

func TestDecode_InvalidVersion(t *testing.T) {
	data := append(buildHeaderBytes(Magic, Version+1, byte(PacketTypeDebugHello), 0, 0, 0, 0), []byte("payload")...)
	_, err := Decode(data)
	if err == nil {
		t.Fatal("Decode(invalid version): expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidVersion) {
		t.Errorf("Decode(invalid version): expected ErrInvalidVersion, got %v", err)
	}
}

func TestDecode_InvalidPacketType(t *testing.T) {
	data := append(buildHeaderBytes(Magic, Version, byte(PacketTypeNone), 0, 0, 0, 0), []byte("payload")...)
	_, err := Decode(data)
	if err == nil {
		t.Fatal("Decode(invalid packet type): expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidPacketType) {
		t.Errorf("Decode(invalid packet type): expected ErrInvalidPacketType, got %v", err)
	}
}

func TestDecode_Valid(t *testing.T) {
	magic := uint8(Magic)
	version := uint8(Version)
	packetType := PacketTypeDebugHello
	flags := uint8(0x42)
	id := uint32(1)
	length := uint32(100)
	fragmentID := uint16(2)
	payload := []byte("payload")
	data := append(buildHeaderBytes(magic, version, byte(packetType), flags, id, length, fragmentID), payload...)
	pkt, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode(valid): unexpected error: %v", err)
	}
	h := pkt.PacketHeader
	if h.Magic != magic {
		t.Errorf("header.Magic = %v, want %v", h.Magic, magic)
	}
	if h.Version != version {
		t.Errorf("header.Version = %v, want %v", h.Version, version)
	}
	if h.PacketType != packetType {
		t.Errorf("header.PacketType = %v, want %v", h.PacketType, packetType)
	}
	if h.Flags != flags {
		t.Errorf("header.Flags = %v, want %v", h.Flags, flags)
	}
	if h.ID != id {
		t.Errorf("header.ID = %v, want %v", h.ID, id)
	}
	if h.Length != length {
		t.Errorf("header.Length = %v, want %v", h.Length, length)
	}
	if h.FragmentID != fragmentID {
		t.Errorf("header.FragmentID = %v, want %v", h.FragmentID, fragmentID)
	}
	if !bytes.Equal(pkt.Payload, payload) {
		t.Errorf("payload = %v, want %v", pkt.Payload, payload)
	}
}

func TestEncodeDecode_RoundTrip(t *testing.T) {
	original := Packet{
		PacketHeader: Header{
			Magic:      Magic,
			Version:    Version,
			PacketType: PacketTypeDebugHello,
			Flags:      0x42,
			ID:         1,
			Length:     100,
			FragmentID: 2,
		},
		Payload: []byte("roundtrip payload"),
	}
	encoded := original.Encode()
	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("Decode(encoded): unexpected error: %v", err)
	}
	h := decoded.PacketHeader
	orig := original.PacketHeader
	if h.Magic != orig.Magic {
		t.Errorf("header.Magic = %v, want %v", h.Magic, orig.Magic)
	}
	if h.Version != orig.Version {
		t.Errorf("header.Version = %v, want %v", h.Version, orig.Version)
	}
	if h.PacketType != orig.PacketType {
		t.Errorf("header.PacketType = %v, want %v", h.PacketType, orig.PacketType)
	}
	if h.Flags != orig.Flags {
		t.Errorf("header.Flags = %v, want %v", h.Flags, orig.Flags)
	}
	if h.ID != orig.ID {
		t.Errorf("header.ID = %v, want %v", h.ID, orig.ID)
	}
	if h.Length != orig.Length {
		t.Errorf("header.Length = %v, want %v", h.Length, orig.Length)
	}
	if h.FragmentID != orig.FragmentID {
		t.Errorf("header.FragmentID = %v, want %v", h.FragmentID, orig.FragmentID)
	}
	if !bytes.Equal(decoded.Payload, original.Payload) {
		t.Errorf("payload = %v, want %v", decoded.Payload, original.Payload)
	}
}

func TestEncodeHeader(t *testing.T) {
	header := Header{
		Magic:      Magic,
		Version:    Version,
		PacketType: PacketTypeDebugHello,
		Flags:      0,
		ID:         1,
		Length:     10,
		FragmentID: 0,
	}
	encoded := EncodeHeader(header)
	if len(encoded) != 14 {
		t.Errorf("EncodeHeader: len(encoded) = %d, want 7", len(encoded))
	}
}

func TestPacket_Encode(t *testing.T) {
	header := Header{
		Magic:      Magic,
		Version:    Version,
		PacketType: PacketTypeDebugHello,
		Flags:      0,
		ID:         0,
		Length:     0,
		FragmentID: 0,
	}
	payload := []byte("payload")
	p := Packet{PacketHeader: header, Payload: payload}
	encoded := p.Encode()
	wantLen := HeaderSize + len(payload)
	if len(encoded) != wantLen {
		t.Errorf("Encode: len(encoded) = %d, want %d", len(encoded), wantLen)
	}
	headerBytes := EncodeHeader(header)
	if len(headerBytes) <= len(encoded) && len(headerBytes) > 0 {
		for i := 0; i < len(headerBytes); i++ {
			if encoded[i] != headerBytes[i] {
				t.Errorf("Encode: encoded[%d] = %v, want %v (from EncodeHeader)", i, encoded[i], headerBytes[i])
				break
			}
		}
	}
}
