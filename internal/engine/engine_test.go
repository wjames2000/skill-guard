package engine

import (
	"os"
	"path/filepath"
	"testing"

	pkgtypes "github.com/hpds.cc/skill-guard/pkg/types"
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
