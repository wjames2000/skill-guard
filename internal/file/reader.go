package file

import (
	"bufio"
	"os"
)

const DefaultBufferSize = 64 * 1024

func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path) // #nosec G304
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, DefaultBufferSize), DefaultBufferSize)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
