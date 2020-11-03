package format

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mpkgr/utils"
	"os"
	"unsafe"

	"hash/crc64"

	"github.com/DataDog/zstd"
)

const contentPayloadVersion = 1

// ContentPayload is the payload for package content
type ContentPayload struct {
	Pt      Payload
	Order   []string
	Content map[string]string
}

// Encode encoes the ContentPayload to the file
func (c *ContentPayload) Encode(fp *os.File) {
	startPoint, err := fp.Seek(0, os.SEEK_CUR)
	if err != nil {
		log.Fatal(err)
	}
	c.Pt.Encode(fp)

	fmt.Print(startPoint)

	switch c.Pt.Compression {
	case payloadCompressionNone:
		c.encodeNoCompression(fp)
		break
	case payloadCompressionZstd:
		c.encodeZstdCompression(fp)
		break
	}

	fp.Sync()

	fp.Seek(startPoint, os.SEEK_SET)

	c.Pt.Encode(fp)

	fp.Seek(int64(c.Pt.Length), os.SEEK_CUR)

}

func (c *ContentPayload) encodeNoCompression(fp *os.File) {
	const chunkSize = 16 * 1024 * 1024

	tab := new(crc64.Table)

	crchash := crc64.New(tab)

	var written uint64 = 0

	for _, k := range c.Order {
		v := c.Content[k]
		file, err := os.Open(v)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err = file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		r := bufio.NewReader(file)
		b := make([]byte, chunkSize)
		for {
			n, err := r.Read(b)
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatal(err)
				} else {
					break
				}
			}
			t, _ := fp.Write(b[0:n])
			utils.Assert(t == int(unsafe.Sizeof(b)), "Failed to write content to payload")
			crchash.Write(b[0:n])
			written += uint64(len(b[0:n]))
		}

		crchash.Sum(c.Pt.CRC64)
		c.Pt.Length = written
		c.Pt.Size = written
	}
}

func (c *ContentPayload) encodeZstdCompression(fp *os.File) {

	defer fp.Close()

	const chunkSize = 16 * 1024 * 1024

	var compSize uint64 = 0
	var normSize uint64 = 0

	tab := new(crc64.Table)
	crchash := crc64.New(tab)

	compressed := io.MultiWriter(fp, crchash)
	compressor := zstd.NewWriter(compressed)

	for _, k := range c.Order {
		v := c.Content[k]

		file, err := os.Open(v)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err = file.Close(); err != nil {
				log.Fatal(err)
			}
		}()

		r := bufio.NewReader(file)
		b := make([]byte, chunkSize)
		for {
			n, err := r.Read(b)
			if err != nil {
				if err.Error() != "EOF" {
					log.Fatal(err)
				} else {
					break
				}
			}

			t, err := compressor.Write(b[0:n])
			if err != nil {
				log.Panic(err)
			}

			normSize += uint64(len(b[0:n]))
			compSize += uint64(t)
		}

		if err := compressor.Close(); err != nil {
			log.Fatal(err)
		}
	}

	crchash.Sum(c.Pt.CRC64)
	c.Pt.Length = compSize
	c.Pt.Size = normSize
}

// AddFile adds a file to the payload
func (c *ContentPayload) AddFile(hashID, sourcePath string) {
	_, ok := c.Content[hashID]
	utils.Assert(!ok, "AddFile(): must be a unique hash")
	c.Content[hashID] = sourcePath
	c.Order = append(c.Order, hashID)
	c.Pt.NumRecords++
}

// DefaultContentPayload creates a default for ContentPayload
func DefaultContentPayload() *ContentPayload {
	r := &ContentPayload{
		Pt: Payload{
			Length:         0,
			Size:           0,
			NumRecords:     0,
			PayloadVersion: contentPayloadVersion,
			Type:           payloadTypeContent,
			Compression:    payloadCompressionNone,
		},
		Order:   []string{},
		Content: make(map[string]string),
	}

	return r
}
