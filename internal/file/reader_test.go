package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("line1\nline2\nline3\n"), 0644)

	lines, err := ReadLines(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 3 {
		t.Fatalf("期望 3 行, 得到 %d", len(lines))
	}
	if lines[0] != "line1" {
		t.Errorf("第一行期望 'line1', 得到 '%s'", lines[0])
	}
}

func TestReadLines_FileNotExist(t *testing.T) {
	_, err := ReadLines("/nonexistent/file.txt")
	if err == nil {
		t.Error("期望错误，但得到 nil")
	}
}
