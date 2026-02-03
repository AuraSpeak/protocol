package protocol

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func init() {
	// Reduce log output during fuzzing so the fuzzer is not slowed and logs are not flooded.
	logrus.SetLevel(logrus.PanicLevel)
}

// FuzzDecode fuzzes Decode with arbitrary byte slices.
// Valid inputs return a packet; invalid inputs must return an error without panicking.
func FuzzDecode(f *testing.F) {
	// Seed corpus: valid packet (DebugHello + empty payload), valid header only, invalid type.
	f.Add([]byte{0x90})
	f.Add([]byte{0x90, 0x00})
	f.Add([]byte{0x01})
	f.Add([]byte{0x91, 0x01, 0x02})
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = Decode(data)
	})
}

// FuzzDecodeHeader fuzzes DecodeHeader with arbitrary byte slices.
// Valid inputs return a header; invalid inputs must return an error without panicking.
func FuzzDecodeHeader(f *testing.F) {
	// Seed corpus: valid types (0x90, 0x01), invalid type (0x00), empty/short.
	f.Add([]byte{0x90})
	f.Add([]byte{0x00})
	f.Add([]byte{0x01})
	f.Add([]byte{0x91})
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = DecodeHeader(data)
	})
}
