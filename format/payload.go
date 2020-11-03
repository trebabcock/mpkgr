package format

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"unsafe"

	"mpkgr/utils"
)

// Payload is a generic payload
type Payload struct {
	Length         uint64
	Size           uint64
	CRC64          []byte
	NumRecords     uint32
	PayloadVersion uint16
	Type           payloadType
	Compression    payloadCompression
}

// Encode encodes the payload to the file
func (p *Payload) Encode(fp *os.File) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, p)
	if err != nil {
		log.Fatal(err)
	}
	t, _ := fp.WriteAt(buf.Bytes(), int64(unsafe.Sizeof(Header{})))
	utils.Assert(t == int(unsafe.Sizeof(Payload{})), "Failed to write payload to file")
}

// DefaultPayload is the default payload
func DefaultPayload() *Payload {
	p := &Payload{
		Length:         12,
		Size:           10,
		CRC64:          []byte{0, 0, 0, 0, 0, 0, 0, 0},
		NumRecords:     4,
		PayloadVersion: 5,
		Type:           payloadTypeMeta,
		Compression:    payloadCompressionZstd,
	}

	return p
}

type payloadType uint8

const (
	payloadTypeUnknown    payloadType = 0
	payloadTypeMeta       payloadType = 1
	payloadTypeContent    payloadType = 2
	payloadTypeLayout     payloadType = 3
	payloadTypeIndex      payloadType = 4
	payloadTypeAttributes payloadType = 5
)

type payloadCompression uint8

const (
	payloadCompressionUnknown payloadCompression = 0
	payloadCompressionNone    payloadCompression = 1
	payloadCompressionZstd    payloadCompression = 2
	payloadCompressionZlib    payloadCompression = 3
)
