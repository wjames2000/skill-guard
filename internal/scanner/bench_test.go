package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func BenchmarkScan_10Files(b *testing.B) {
	dir := b.TempDir()
	for i := 0; i < 10; i++ {
		os.WriteFile(filepath.Join(dir, fmt.Sprintf("test%d.py", i)), []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0644)
	}

	cfg := &pkgtypes.Config{
		Paths:   []string{dir},
		MaxSize: 10 * 1024 * 1024,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Scan(cfg)
	}
}

func BenchmarkScan_100Files(b *testing.B) {
	dir := b.TempDir()
	for i := 0; i < 100; i++ {
		os.WriteFile(filepath.Join(dir, fmt.Sprintf("test%d.py", i)), []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0644)
	}

	cfg := &pkgtypes.Config{
		Paths:   []string{dir},
		MaxSize: 10 * 1024 * 1024,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Scan(cfg)
	}
}
