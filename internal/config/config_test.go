package config

import (
	"os"
	"path/filepath"
	"testing"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func TestLoadFile_NotExist(t *testing.T) {
	cfg, err := LoadFile("/nonexistent/.skillguard.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if cfg == nil {
		t.Fatal("配置文件不存在应返回空配置")
	}
}

func TestLoadFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".skillguard.yaml")
	os.WriteFile(path, []byte("ignore:\n  - \"tests/**\"\nseverity: \"high\"\nmax-size: \"5MB\"\nconcurrency: 4\n"), 0644)
	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Ignore) != 1 || cfg.Ignore[0] != "tests/**" {
		t.Errorf("Ignore 不匹配: %v", cfg.Ignore)
	}
	if cfg.Severity != "high" {
		t.Errorf("Severity 应为 high, 得到 %s", cfg.Severity)
	}
	if cfg.Concurrency != 4 {
		t.Errorf("Concurrency 应为 4, 得到 %d", cfg.Concurrency)
	}
}

func TestMergeWithCLI(t *testing.T) {
	cli := pkgtypes.DefaultConfig()
	cli.Ignore = []string{"--cli-ignore"}
	fileCfg := &FileConfig{Ignore: []string{"--file-ignore"}, Severity: "high", MaxSize: "5MB"}
	merged := MergeWithCLI(cli, fileCfg)
	if len(merged.Ignore) != 2 {
		t.Errorf("Ignore 长度应为 2, 得到 %d", len(merged.Ignore))
	}
	if merged.Severity != "high" {
		t.Errorf("Severity 应为 high, 得到 %s", merged.Severity)
	}
}

func TestParseMaxSize(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"5MB", 5 * 1024 * 1024},
		{"10MB", 10 * 1024 * 1024},
		{"1GB", 1 * 1024 * 1024 * 1024},
		{"500KB", 500 * 1024},
		{"invalid", 0},
	}
	for _, tt := range tests {
		got := parseMaxSize(tt.input)
		if got != tt.want {
			t.Errorf("parseMaxSize(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestLoadFile_Malformed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".skillguard.yaml")
	os.WriteFile(path, []byte(": invalid yaml [[["), 0644)
	_, err := LoadFile(path)
	if err == nil {
		t.Error("格式错误的 YAML 应返回错误")
	}
}

func TestMergeWithCLI_Override(t *testing.T) {
	cli := &pkgtypes.Config{Severity: "high", MaxSize: 5 * 1024 * 1024}
	fileCfg := &FileConfig{Severity: "low", MaxSize: "1MB"}
	merged := MergeWithCLI(cli, fileCfg)
	if merged.Severity != "high" {
		t.Errorf("CLI severity 应覆盖文件配置: 得到 %s", merged.Severity)
	}
	if merged.MaxSize != 5*1024*1024 {
		t.Errorf("CLI MaxSize 应覆盖文件配置: 得到 %d", merged.MaxSize)
	}
}
