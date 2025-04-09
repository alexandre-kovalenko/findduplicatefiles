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

// /////////////////////////////////////////////////////////////////////////////
// Make new file unit
func MakeFileUnit(filename string) (FileUnit, error) {
	file, err := os.Open(filename)
	if err != nil {
		return FileUnit{}, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return FileUnit{}, err
	}

	return FileUnit{
		Name:     filename,
		Checksum: hash.Sum(nil),
	}, nil
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
