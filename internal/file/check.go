package file

import (
	"os"
	"unicode/utf8"
)

const checkHeaderSize = 512

func IsValidUTF8(path string) bool {
	f, err := os.Open(path) // #nosec G304
	if err != nil {
		return false
	}
	defer f.Close()

	header := make([]byte, checkHeaderSize)
	n, _ := f.Read(header)
	return utf8.Valid(header[:n])
}

func IsWithinSizeLimit(path string, maxSize int64) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Size() <= maxSize
}
