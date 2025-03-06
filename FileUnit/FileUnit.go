package FileUnit

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

type FileUnit struct {
	FileName string
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
		FileName: filename,
		Checksum: hash.Sum(nil),
	}, nil
}

// ///////////////////////////////////////////////////////////////////////////////
// Make string representation of the file unit
func (f FileUnit) ToString() string {
	result := f.FileName + ": 0x"
	for _, b := range f.Checksum {
		result += fmt.Sprintf("%02x", int(b)&0xFF)
	}
	return result
}
