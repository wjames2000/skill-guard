package engine

import (
	"os"
	"path/filepath"
	"testing"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func BenchmarkEngine_Match_SingleRule(b *testing.B) {
	eng, _ := New("", false)
	dir := b.TempDir()
	path := filepath.Join(dir, "test.py")
	os.WriteFile(path, []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "test.py", Ext: ".py"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eng.Match(target)
	}
}

func BenchmarkEngine_Match_AllRules(b *testing.B) {
	eng, _ := New("", false)
	dir := b.TempDir()
	path := filepath.Join(dir, "mixed.py")
	os.WriteFile(path, []byte(`
key = "AKIAIOSFODNN7EXAMPLE"
os.system("ls -la")
chmod 777 /var/www
eval(input())
import subprocess
subprocess.call(["ls"])
`), 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "mixed.py", Ext: ".py"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eng.Match(target)
	}
}

func BenchmarkEngine_Match_All100Rules(b *testing.B) {
	eng, _ := New("", false)
	dir := b.TempDir()
	content := []byte{}
	for i := 0; i < 100; i++ {
		content = append(content, []byte("key = \"AKIAIOSFODNN7EXAMPLE\"\nos.system(\"ls\")\nchmod 777 /var\ncat ~/.ssh/id_rsa\n")...)
	}
	path := filepath.Join(dir, "mixed.py")
	os.WriteFile(path, content, 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "mixed.py", Ext: ".py"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eng.Match(target)
	}
}

func BenchmarkEngine_Match_EmptyFile(b *testing.B) {
	eng, _ := New("", false)
	dir := b.TempDir()
	path := filepath.Join(dir, "empty.py")
	os.WriteFile(path, []byte{}, 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "empty.py", Ext: ".py"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eng.Match(target)
	}
}
