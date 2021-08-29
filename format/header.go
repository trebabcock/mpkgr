package format

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"

	"mpkgr/utils"
)

const mpkgFileHeader uint32 = 0x6d706b67

func integrityCheck() [21]byte {
	return [21]byte{2, 6, 5, 4, 7, 3, 3, 2, 4, 7, 5, 1, 2, 3, 4, 7, 4, 6, 2, 7, 1}
}

const mpkgFormatVersionNumber = 1

// MpkgFileType is the file type for the header
type MpkgFileType uint8

const (
	mpkgFileTypeUnkown MpkgFileType = 0
	mpkgFileTypeBinary MpkgFileType = 1
	mpkgFileTypeDelta  MpkgFileType = 2
)

// Header is the file header
type Header struct {
	Magic         uint32
	NumPayloads   uint16
	Padding       [21]byte
	Type          MpkgFileType
	VersionNumber uint32
}

// ConstructHeader is a constructor for the header
func ConstructHeader(versionNumber uint32) *Header {
	h := &Header{
		Magic:         mpkgFileHeader,
		NumPayloads:   0,
		Padding:       integrityCheck(),
		Type:          mpkgFileTypeBinary,
		VersionNumber: versionNumber,
	}

	return h
}

// Validate valides the header
func (h *Header) Validate() {
	if h.Magic != mpkgFileHeader {
		log.Fatal("Header.Validate: invalid header")
	}

	if h.Padding != integrityCheck() {
		log.Fatal("Header.Validate: corrupt integrity")
	}

	if h.Type == mpkgFileTypeUnkown {
		log.Fatal("Header.Validate: unknown package type")
	}

	if h.VersionNumber != mpkgFormatVersionNumber {
		log.Fatal("Header.Validate: unsupported package version")
	}
}

// Encode encodes the header to the file
func (h *Header) Encode(fp *os.File) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, h)
	if err != nil {
		log.Fatal(err)
	}
	t, _ := fp.Write(buf.Bytes())
	utils.Assert(t == int(unsafe.Sizeof(Header{})), "Failed to write header to file")
}

// CheckHeaderSize makes sure the header size is 32 bytes
func CheckHeaderSize() {
	utils.Assert(int(unsafe.Sizeof(Header{})) == 32, "Header must be 32-bytes, found: "+fmt.Sprint(unsafe.Sizeof(Header{}))+" bytes")
}
