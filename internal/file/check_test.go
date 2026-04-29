package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidUTF8(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "utf8.txt")
	os.WriteFile(path, []byte("hello 世界"), 0644)
	if !IsValidUTF8(path) {
		t.Error("UTF-8 文件应返回 true")
	}
}

func TestIsWithinSizeLimit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello"), 0644)

	if !IsWithinSizeLimit(path, 100) {
		t.Error("小文件应在限制内")
	}
	if IsWithinSizeLimit(path, 1) {
		t.Error("超限文件应返回 false")
	}
}
