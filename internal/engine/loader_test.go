package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRulesFile_YAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.yaml")
	os.WriteFile(path, []byte("rules:\n  - id: \"USR-001\"\n    name: \"自定义\"\n    severity: \"High\"\n    pattern: \"malicious\"\n    keywords: [\"malicious\"]\n"), 0644)
	rules, err := LoadRulesFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("期望 1 条规则, 得到 %d", len(rules))
	}
	if rules[0].ID != "USR-001" {
		t.Errorf("ID 应为 USR-001, 得到 %s", rules[0].ID)
	}
}

func TestLoadRulesFile_JSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.json")
	os.WriteFile(path, []byte(`[{"id":"USR-001","name":"Test","severity":"High","pattern":"test"}]`), 0644)
	rules, err := LoadRulesFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("期望 1 条规则, 得到 %d", len(rules))
	}
}

func TestLoadRulesFile_InvalidExt(t *testing.T) {
	_, err := LoadRulesFile("rules.txt")
	if err == nil {
		t.Error("不支持的文件格式应返回错误")
	}
}
