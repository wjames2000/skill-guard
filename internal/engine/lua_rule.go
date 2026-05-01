package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

// LuaRule 表示一个 LUA 脚本规则
type LuaRule struct {
	ID          string
	Name        string
	Severity    string
	Description string
	FileTypes   []string
	scriptFile  string
	sandbox     *LuaSandbox
}

// LoadLuaRule 从 .lua 文件加载规则
// 文件必须返回包含 id/name/severity/description/file_types 的 rule 表，
// 并提供 rule.match(line, filename) 匹配函数。
func LoadLuaRule(path string) (*LuaRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 LUA 规则失败: %w", err)
	}

	tempSandbox := NewLuaSandbox()
	defer tempSandbox.Close()

	// 执行脚本以注册 rule 全局表
	fn, err := tempSandbox.L.LoadString(string(data))
	if err != nil {
		return nil, fmt.Errorf("加载 LUA 规则失败: %w", err)
	}

	// Push and call the loaded function
	tempSandbox.L.Push(fn)
	if err := tempSandbox.L.PCall(0, lua.MultRet, nil); err != nil {
		return nil, fmt.Errorf("执行 LUA 规则失败: %w", err)
	}

	// Pop the return value(s) - we don't need them, we use globals
	// (PCall with MultRet pushes all return values; clean them up)
	n := tempSandbox.L.GetTop()
	for i := 0; i < n; i++ {
		tempSandbox.L.Pop(1)
	}

	// Extract rule table from global
	ruleTable := tempSandbox.L.GetGlobal("rule")
	if ruleTable.Type() != lua.LTTable {
		return nil, fmt.Errorf("LUA 规则必须定义全局 'rule' 表（文件: %s）", path)
	}

	// Create a dedicated sandbox for the rule
	rule := &LuaRule{
		scriptFile: path,
		sandbox:    NewLuaSandbox(),
	}

	if v := tempSandbox.L.GetField(ruleTable, "id"); v.Type() == lua.LTString {
		rule.ID = v.String()
	}
	if v := tempSandbox.L.GetField(ruleTable, "name"); v.Type() == lua.LTString {
		rule.Name = v.String()
	}
	if v := tempSandbox.L.GetField(ruleTable, "severity"); v.Type() == lua.LTString {
		rule.Severity = v.String()
	}
	if v := tempSandbox.L.GetField(ruleTable, "description"); v.Type() == lua.LTString {
		rule.Description = v.String()
	}
	if v := tempSandbox.L.GetField(ruleTable, "file_types"); v.Type() == lua.LTTable {
		var types []string
		v.(*lua.LTable).ForEach(func(_, val lua.LValue) {
			if val.Type() == lua.LTString {
				types = append(types, val.String())
			}
		})
		rule.FileTypes = types
	}

	if rule.ID == "" {
		return nil, fmt.Errorf("LUA 规则缺少必填字段 'id'（文件: %s）", path)
	}

	return rule, nil
}

// MatchesFileType 检查规则是否适用于指定文件扩展名
func (r *LuaRule) MatchesFileType(ext string) bool {
	if len(r.FileTypes) == 0 {
		return true
	}
	ext = strings.ToLower(ext)
	for _, ft := range r.FileTypes {
		if strings.EqualFold(ext, ft) {
			return true
		}
	}
	return false
}

// MatchLine 对单行内容执行 LUA 规则的 match 函数
// 每次调用重新加载脚本以确保沙箱隔离
func (r *LuaRule) MatchLine(line string, filename string) *pkgtypes.MatchResult {
	// 每次匹配重新加载脚本，避免状态污染
	fn, err := r.sandbox.L.LoadFile(r.scriptFile)
	if err != nil {
		return nil
	}
	r.sandbox.L.Push(fn)
	if err := r.sandbox.L.PCall(0, lua.MultRet, nil); err != nil {
		return nil
	}
	// 清理返回的栈（脚本的 return 值）
	n := r.sandbox.L.GetTop()
	for i := 0; i < n; i++ {
		r.sandbox.L.Pop(1)
	}

	// 获取 rule 表
	ruleTable := r.sandbox.L.GetGlobal("rule")
	if ruleTable.Type() != lua.LTTable {
		return nil
	}

	// 获取 match 函数
	matchFn := r.sandbox.L.GetField(ruleTable, "match")
	if matchFn.Type() != lua.LTFunction {
		return nil
	}

	// 调用 match(line, filename)
	if err := r.sandbox.L.CallByParam(lua.P{
		Fn:      matchFn,
		NRet:    2,
		Protect: true,
	}, lua.LString(line), lua.LString(filename)); err != nil {
		return nil
	}

	matched := r.sandbox.L.Get(-2)
	reason := r.sandbox.L.Get(-1)
	r.sandbox.L.Pop(2)

	if matched == lua.LTrue {
		content := reason.String()
		if len(content) > 120 {
			content = content[:120] + "..."
		}
		return &pkgtypes.MatchResult{
			RuleID:      r.ID,
			RuleName:    r.Name,
			Severity:    r.Severity,
			LineContent: content,
			MatchType:   "lua",
		}
	}
	return nil
}

// LoadLuaRulesFromDir 从指定目录加载所有 .lua 规则文件
func LoadLuaRulesFromDir(luaDir string) ([]*LuaRule, error) {
	entries, err := os.ReadDir(luaDir)
	if err != nil {
		return nil, fmt.Errorf("读取 LUA 规则目录失败: %w", err)
	}

	var rules []*LuaRule
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".lua") {
			continue
		}
		rule, err := LoadLuaRule(filepath.Join(luaDir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("加载 LUA 规则 %s 失败: %w", entry.Name(), err)
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
