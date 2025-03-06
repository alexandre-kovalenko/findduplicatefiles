package FileUnit

type FileUnit struct {
	FileName string
	Checksum []byte
}

func MakeFileUnit(filename path) (FileUnit, error) {
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