package format

import (
	"mpkgr/utils"
	"os"
)

// MpkgWriter writes the mpkg file
type MpkgWriter struct {
	File         *os.File
	FileHeader   Header
	FilePayloads []*Payload
}

// FileType returns the file type
func (m *MpkgWriter) FileType() MpkgFileType {
	return m.FileHeader.Type
}

// SetFileType sets the file type
func (m *MpkgWriter) SetFileType(fileType MpkgFileType) {
	m.FileHeader.Type = fileType
}

// Close closes the file
func (m *MpkgWriter) Close() error {
	m.File.Seek(0, 0)
	m.FileHeader.Encode(m.File)
	err := m.File.Close()
	if err != nil {
		return err
	}
	return nil
}

// AddPayload adds a payload to the writer
func (m *MpkgWriter) AddPayload(payload *Payload) {
	m.FilePayloads = append(m.FilePayloads, payload)
	m.FileHeader.NumPayloads++
}

// Flush flushes the writer buffer
func (m *MpkgWriter) Flush() {
	m.File.Seek(0, 0)
	m.FileHeader.Encode(m.File)

	for _, p := range m.FilePayloads {
		switch p.Type {
		case payloadTypeMeta:

		case payloadTypeContent:
		case payloadTypeLayout:
		case payloadTypeIndex:
		default:
			utils.Assert(false, "Unsupported type: "+string(p.Type))
		}
	}
}

// NewMpkgWriter creates a new MpkgWriter
func NewMpkgWriter(file *os.File, versionNumber uint32) *MpkgWriter {
	w := &MpkgWriter{
		File:       file,
		FileHeader: *ConstructHeader(versionNumber),
	}

	w.FileHeader.NumPayloads = 0
	w.FileHeader.Encode(file)

	return w
}
