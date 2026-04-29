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

func TestDiscover_WithSymlink(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "real.py"), []byte("print('safe')"), 0644)
	err := os.Symlink(filepath.Join(dir, "real.py"), filepath.Join(dir, "link.py"))
	if err != nil {
		t.Skip("不支持符号链接:", err)
	}
	files, _ := Discover([]string{dir}, nil)
	for _, f := range files {
		if f.RelPath == "link.py" {
			t.Error("符号链接应被跳过")
		}
	}
}

func TestDiscover_LargeFile(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "small.py"), []byte("x=1"), 0644)
	largeData := make([]byte, 100)
	os.WriteFile(filepath.Join(dir, "large.py"), largeData, 0644)

	files, _ := Discover([]string{dir}, &DiscoverOpts{MaxSize: 50})
	if len(files) != 1 {
		t.Errorf("期望 1 个文件(小文件), 得到 %d", len(files))
	}
	if len(files) > 0 && files[0].RelPath != "small.py" {
		t.Errorf("应保留 small.py, 得到 %s", files[0].RelPath)
	}
}

func TestDiscover_NonUTF8(t *testing.T) {
	dir := t.TempDir()
	// .py is a default extension; non-UTF8 content should not crash discover
	path := filepath.Join(dir, "binary.py")
	os.WriteFile(path, []byte{0xff, 0xfe, 0x00, 0x01}, 0644)
	files, _ := Discover([]string{dir}, nil)
	found := false
	for _, f := range files {
		if f.RelPath == "binary.py" {
			found = true
			break
		}
	}
	if !found {
		t.Error("binary.py 应被 discover 发现（engine 层会跳过）")
	}
}

func TestDiscover_EmptyRoot(t *testing.T) {
	files, err := Discover([]string{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Errorf("空路径应返回空，得到 %d", len(files))
	}
}

func TestDiscover_IgnorePattern(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.py"), []byte("x=1"), 0644)
	os.WriteFile(filepath.Join(dir, "test.md"), []byte("# doc"), 0644)
	files, _ := Discover([]string{dir}, &DiscoverOpts{ExtExclude: []string{".md"}})
	for _, f := range files {
		if f.Ext == ".md" {
			t.Error(".md 文件应被排除")
		}
	}
}
