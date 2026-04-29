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

func TestDiscover_WithGitIgnore(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("ignored_dir/\n*.log\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "ignored_dir"), 0755)
	os.WriteFile(filepath.Join(dir, "ignored_dir", "secret.py"), []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0644)
	os.WriteFile(filepath.Join(dir, "keep.py"), []byte("print('ok')"), 0644)
	os.WriteFile(filepath.Join(dir, "debug.log"), []byte("log content"), 0644)

	files, err := Discover([]string{dir}, &DiscoverOpts{
		DiscoverGitIgnore: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		t.Logf("found: %s", f.RelPath)
	}
	if len(files) != 1 {
		t.Errorf("期望 1 个文件 (keep.py), 得到 %d", len(files))
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
