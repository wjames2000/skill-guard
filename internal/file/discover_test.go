package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscover(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.py"), []byte("print('hello')"), 0644)
	os.WriteFile(filepath.Join(dir, "test.md"), []byte("# readme"), 0644)
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)
	os.WriteFile(filepath.Join(dir, ".git", "config"), []byte(""), 0644)
	os.MkdirAll(filepath.Join(dir, "subdir"), 0755)
	os.WriteFile(filepath.Join(dir, "subdir", "script.sh"), []byte("echo hi"), 0644)

	files, err := Discover([]string{dir}, &DiscoverOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 {
		t.Errorf("期望 3 个文件, 得到 %d", len(files))
	}
	for _, f := range files {
		if f.RelPath == ".git/config" || strings.HasPrefix(f.RelPath, ".git/") {
			t.Errorf(".git 目录应被忽略: %s", f.RelPath)
		}
	}
}

func TestDiscover_ExtInclude(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.py"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "test.md"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "test.sh"), []byte(""), 0644)

	files, _ := Discover([]string{dir}, &DiscoverOpts{
		ExtInclude: []string{".py", ".sh"},
	})
	if len(files) != 2 {
		t.Errorf("期望 2 个文件(.py/.sh), 得到 %d", len(files))
	}
}
