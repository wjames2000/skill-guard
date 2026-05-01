package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLuaRule(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.lua")
	content := `rule = {id="TEST-001",name="Test",severity="High",description="Test rule"}
function rule.match(line, filename)
    if line:find("danger") then return true, "danger found" end
    return false, ""
end
return rule`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	rule, err := LoadLuaRule(path)
	if err != nil {
		t.Fatal(err)
	}
	if rule.ID != "TEST-001" {
		t.Errorf("ID 期望 TEST-001, 得到 %s", rule.ID)
	}
	if rule.Severity != "High" {
		t.Errorf("Severity 期望 High, 得到 %s", rule.Severity)
	}
	if rule.Name != "Test" {
		t.Errorf("Name 期望 Test, 得到 %s", rule.Name)
	}
}

func TestLoadLuaRule_MissingTable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.lua")
	if err := os.WriteFile(path, []byte("x = 1"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadLuaRule(path)
	if err == nil {
		t.Error("缺少 rule 表应返回错误")
	}
}

func TestLuaRule_MatchLine_Found(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.lua")
	content := `rule = {id="TEST-002",name="Match Test",severity="High",description="Test"}
function rule.match(line, filename)
    if line:find("secret") then return true, "found secret: " .. line end
    return false, ""
end
return rule`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	rule, err := LoadLuaRule(path)
	if err != nil {
		t.Fatal(err)
	}

	result := rule.MatchLine(`password = "secret123"`, "test.py")
	if result == nil {
		t.Fatal("应匹配到 secret")
	}
	if result.RuleID != "TEST-002" {
		t.Errorf("RuleID 期望 TEST-002, 得到 %s", result.RuleID)
	}
	if result.MatchType != "lua" {
		t.Errorf("MatchType 期望 lua, 得到 %s", result.MatchType)
	}
}

func TestLuaRule_MatchLine_NotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.lua")
	content := `rule = {id="TEST-003",name="No Match",severity="Low",description="Test"}
function rule.match(line, filename)
    if line:find("danger") then return true, "found" end
    return false, ""
end
return rule`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	rule, err := LoadLuaRule(path)
	if err != nil {
		t.Fatal(err)
	}

	result := rule.MatchLine(`safe line of code`, "test.py")
	if result != nil {
		t.Error("安全行不应返回匹配结果")
	}
}

func TestLuaRule_MatchesFileType(t *testing.T) {
	rule := &LuaRule{
		ID:        "TEST-004",
		FileTypes: []string{".py", ".sh"},
	}

	tests := []struct {
		ext  string
		want bool
	}{
		{".py", true},
		{".sh", true},
		{".PY", true},
		{".SH", true},
		{".md", false},
		{".js", false},
	}
	for _, tt := range tests {
		if got := rule.MatchesFileType(tt.ext); got != tt.want {
			t.Errorf("MatchesFileType(%q) = %v, want %v", tt.ext, got, tt.want)
		}
	}
}

func TestLuaRule_MatchesFileType_Empty(t *testing.T) {
	rule := &LuaRule{ID: "TEST-005"}
	if !rule.MatchesFileType(".anything") {
		t.Error("未限定 FileTypes 应匹配所有扩展名")
	}
}

func TestLoadLuaRule_MissingID(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing_id.lua")
	content := `rule = {name="No ID",severity="High"}
function rule.match(line, filename) return false, "" end
return rule`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadLuaRule(path)
	if err == nil {
		t.Error("缺少 id 字段应返回错误")
	}
}

func TestLoadLuaRulesFromDir(t *testing.T) {
	dir := t.TempDir()

	// Create a valid rule
	rule1 := `rule = {id="DIR-001",name="Dir Rule 1",severity="Low",description="Test"}
function rule.match(line, filename) return false, "" end
return rule`
	if err := os.WriteFile(filepath.Join(dir, "rule1.lua"), []byte(rule1), 0644); err != nil {
		t.Fatal(err)
	}

	// Create another valid rule
	rule2 := `rule = {id="DIR-002",name="Dir Rule 2",severity="Medium",description="Test"}
function rule.match(line, filename) return false, "" end
return rule`
	if err := os.WriteFile(filepath.Join(dir, "rule2.lua"), []byte(rule2), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a non-lua file (should be ignored)
	if err := os.WriteFile(filepath.Join(dir, "note.txt"), []byte("not a rule"), 0644); err != nil {
		t.Fatal(err)
	}

	rules, err := LoadLuaRulesFromDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 2 {
		t.Errorf("期望 2 个规则，得到 %d", len(rules))
	}
}

func TestLuaRule_MatchLine_VerifySandboxIsolation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sandbox.lua")
	content := `rule = {id="SANDBOX-001",name="Sandbox Test",severity="High",description="Test"}
function rule.match(line, filename)
    -- 验证 os 和 io 已被移除
    if os ~= nil then return true, "os 未被移除" end
    if io ~= nil then return true, "io 未被移除" end
    return false, ""
end
return rule`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	rule, err := LoadLuaRule(path)
	if err != nil {
		t.Fatal(err)
	}

	// The match function checks that os and io are nil
	result := rule.MatchLine(`test line`, "test.py")
	if result != nil {
		t.Errorf("沙箱检查失败: %s", result.LineContent)
	}
}
