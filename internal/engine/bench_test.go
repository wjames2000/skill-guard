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
