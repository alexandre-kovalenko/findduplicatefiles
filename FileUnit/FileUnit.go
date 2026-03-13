package FileUnit

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

type FileUnit struct {
	Name     string
	Checksum []byte
}

type MakeFileUnitResult struct {
	FileUnit FileUnit
	Error    error
}

const (
	// Size of the buffer to use with io.Copybuffer
	IO_BUFFER_SIZE = 1024 * 1024
)

// /////////////////////////////////////////////////////////////////////////////
// Make new file unit
func MakeFileUnit(filename string, ch chan MakeFileUnitResult) {
	file, err := os.Open(filename)
	if err != nil {
		ch <- MakeFileUnitResult{
			FileUnit: FileUnit{},
			Error:    err,
		}
	}
	defer file.Close()

	hash := sha256.New()
	buffer := make([]byte, IO_BUFFER_SIZE)
	if _, err := io.CopyBuffer(hash, file, buffer); err != nil {
		ch <- MakeFileUnitResult{
			FileUnit: FileUnit{},
			Error:    err,
		}
	}

	ch <- MakeFileUnitResult{
		FileUnit: FileUnit{
			Name:     filename,
			Checksum: hash.Sum(nil),
		},
		Error: nil,
	}
}

// ///////////////////////////////////////////////////////////////////////////////
// Make string representation of the file unit
func (f FileUnit) ToString() string {
	result := f.Name + ": 0x"
	for _, b := range f.Checksum {
		result += fmt.Sprintf("%02x", int(b)&0xFF)
	}
	return result
}

// /////////////////////////////////////////////////////////////////////////////////
// Get base 64 encoded checksum for representation purposes
func (f FileUnit) GetEncodedChecksum() string {
	return base64.URLEncoding.EncodeToString(f.Checksum)
}
