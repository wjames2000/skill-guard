package engine

import (
	"os"
	"path/filepath"
	"testing"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func TestEngine_Match_FindsRisk(t *testing.T) {
	eng, err := New("", false)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "test.py")
	os.WriteFile(path, []byte(`access_key = "AKIAIOSFODNN7EXAMPLE"`), 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "test.py", Ext: ".py"}
	results := eng.Match(target)
	if len(results) == 0 {
		t.Fatal("应检测到风险")
	}
	if results[0].RuleID != "SKL-001" {
		t.Errorf("期望 SKL-001, 得到 %s", results[0].RuleID)
	}
}

func TestEngine_Match_CleanFile(t *testing.T) {
	eng, _ := New("", false)
	dir := t.TempDir()
	path := filepath.Join(dir, "safe.py")
	os.WriteFile(path, []byte(`print("hello world")`), 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "safe.py", Ext: ".py"}
	results := eng.Match(target)
	if len(results) != 0 {
		t.Errorf("安全文件不应有匹配, 得到 %d", len(results))
	}
}

func TestEngine_Match_BinaryFile(t *testing.T) {
	eng, _ := New("", false)
	dir := t.TempDir()
	path := filepath.Join(dir, "binary.bin")
	os.WriteFile(path, []byte{0xff, 0xfe, 0x00, 0x01}, 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "binary.bin", Ext: ".bin"}
	results := eng.Match(target)
	if len(results) != 0 {
		t.Errorf("二进制文件应跳过, 得到 %d 个结果", len(results))
	}
}

func TestEngine_Match_EmptyFile(t *testing.T) {
	eng, _ := New("", false)
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.py")
	os.WriteFile(path, []byte{}, 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "empty.py", Ext: ".py"}
	results := eng.Match(target)
	if len(results) != 0 {
		t.Errorf("空文件应无结果, 得到 %d 个结果", len(results))
	}
}

func TestEngine_Match_MultipleRisks(t *testing.T) {
	eng, _ := New("", false)
	dir := t.TempDir()
	content := []byte("key = \"AKIAIOSFODNN7EXAMPLE\"\nos.system(\"ls\")\n")
	path := filepath.Join(dir, "risk.py")
	os.WriteFile(path, content, 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "risk.py", Ext: ".py"}
	results := eng.Match(target)
	if len(results) < 2 {
		t.Errorf("风险文件应检测到 ≥2 个结果, 得到 %d", len(results))
	}
}

func TestEngine_Match_KeywordsFilter(t *testing.T) {
	rule, _ := NewRule(&pkgtypes.Rule{
		ID: "SKL-011", Pattern: `(mysql|postgres)://[^:]+:[^@]+@`,
		Keywords: []string{"mysql://"},
		Severity: "High",
	})
	if rule.HasKeywords([]string{"no keywords here"}) {
		t.Error("无关键字文件应返回 false")
	}
	if !rule.HasKeywords([]string{"mysql://user:pass@localhost"}) {
		t.Error("含关键字文件应返回 true")
	}
}

func TestEngine_Match_FileTypeLimit(t *testing.T) {
	rule, _ := NewRule(&pkgtypes.Rule{
		ID: "SKL-003", Pattern: `os\.system\s*\(`, Severity: "High",
		FileTypes: []string{".py"},
	})
	if rule.MatchesFileType(".go") {
		t.Error(".go 文件不应匹配 .py 限定规则")
	}
	if !rule.MatchesFileType(".py") {
		t.Error(".py 文件应匹配 .py 限定规则")
	}
}

func TestEngine_Match_BinaryFileEncoded(t *testing.T) {
	eng, _ := New("", false)
	dir := t.TempDir()
	path := filepath.Join(dir, "mixed.txt")
	content := make([]byte, 1000)
	for i := range content {
		content[i] = byte(32 + i%95)
	}
	os.WriteFile(path, content, 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "mixed.txt", Ext: ".txt"}
	results := eng.Match(target)
	_ = results
}

func TestEngine_Match_NoRules(t *testing.T) {
	eng, err := New("", true)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "test.py")
	os.WriteFile(path, []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0644)
	target := &pkgtypes.FileTarget{Path: path, RelPath: "test.py", Ext: ".py"}
	results := eng.Match(target)
	if len(results) != 0 {
		t.Errorf("禁用内置规则应无匹配，得到 %d", len(results))
	}
}
