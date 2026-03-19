package FileUnit

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"rabbitsden.online/FindDuplicateFiles/UnitLimiter"
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
func MakeFileUnit(filename string, ch chan MakeFileUnitResult, l *UnitLimiter.Limiter) {
	l.Acquire()
	// log.Printf("Acquired counter for file '%s'\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		l.Release()
		// log.Printf("Released counter for file '%s'\n", filename)
		ch <- MakeFileUnitResult{
			FileUnit: FileUnit{},
			Error:    err,
		}
	}
	// log.Printf("Opened file '%s' for hashing\n", filename)
	hash := sha256.New()
	buffer := make([]byte, IO_BUFFER_SIZE)
	if _, err := io.CopyBuffer(hash, file, buffer); err != nil {
		file.Close()
		l.Release()
		// log.Printf("Released counter for file '%s'\n", filename)
		ch <- MakeFileUnitResult{
			FileUnit: FileUnit{},
			Error:    err,
		}
	}
	file.Close()
	// log.Printf("Hashed file '%s'\n", filename)
	l.Release()
	// log.Printf("Released counter for file '%s'\n", filename)
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
