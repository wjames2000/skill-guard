package engine

import (
	"fmt"

	"github.com/wjames2000/skill-guard/internal/file"
	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

type Engine struct {
	rules    []*Rule
	luaRules []*LuaRule
}

func New(rulesFile string, disableBuiltin bool) (*Engine, error) {
	eng := &Engine{}
	if !disableBuiltin {
		for _, r := range BuiltinRules() {
			rule, err := NewRule(r)
			if err != nil {
				return nil, err
			}
			eng.rules = append(eng.rules, rule)
		}
	}
	if rulesFile != "" {
		customRules, err := LoadRulesFile(rulesFile)
		if err != nil {
			return nil, err
		}
		if err := eng.mergeCustom(customRules); err != nil {
			return nil, err
		}
	}
	return eng, nil
}

// Match 对目标文件执行所有规则匹配（正则 + LUA）
func (e *Engine) Match(target *pkgtypes.FileTarget) []*pkgtypes.MatchResult {
	results := e.matchRegex(target)
	results = append(results, e.matchLua(target)...)
	return results
}

// matchRegex 执行正则规则匹配
func (e *Engine) matchRegex(target *pkgtypes.FileTarget) []*pkgtypes.MatchResult {
	if !file.IsValidUTF8(target.Path) {
		return nil
	}
	if !file.IsWithinSizeLimit(target.Path, 10*1024*1024) {
		return nil
	}
	lines, err := file.ReadLines(target.Path)
	if err != nil {
		return nil
	}
	var results []*pkgtypes.MatchResult
	for _, rule := range e.rules {
		if !rule.MatchesFileType(target.Ext) {
			continue
		}
		if !rule.HasKeywords(lines) {
			continue
		}
		for lineNum, line := range lines {
			result := rule.MatchLine(line, target.RelPath, lineNum+1)
			if result != nil {
				results = append(results, result)
				break
			}
		}
	}
	return results
}

// LoadLuaRules 从指定目录加载所有 .lua 规则文件
func (e *Engine) LoadLuaRules(luaDir string) error {
	rules, err := LoadLuaRulesFromDir(luaDir)
	if err != nil {
		return err
	}
	e.luaRules = append(e.luaRules, rules...)
	return nil
}

// matchLua 执行 LUA 脚本规则匹配
func (e *Engine) matchLua(target *pkgtypes.FileTarget) []*pkgtypes.MatchResult {
	if len(e.luaRules) == 0 {
		return nil
	}
	if !file.IsValidUTF8(target.Path) {
		return nil
	}
	if !file.IsWithinSizeLimit(target.Path, 10*1024*1024) {
		return nil
	}
	lines, err := file.ReadLines(target.Path)
	if err != nil {
		return nil
	}
	var results []*pkgtypes.MatchResult
	for _, luaRule := range e.luaRules {
		if !luaRule.MatchesFileType(target.Ext) {
			continue
		}
		for lineNum, line := range lines {
			result := luaRule.MatchLine(line, target.RelPath)
			if result != nil {
				result.FilePath = target.RelPath
				result.LineNumber = lineNum + 1
				results = append(results, result)
				break
			}
		}
	}
	return results
}

func (e *Engine) mergeCustom(custom []*pkgtypes.Rule) error {
	existing := make(map[string]int)
	for i, r := range e.rules {
		existing[r.ID] = i
	}
	for _, cr := range custom {
		if cr.ID == "" || cr.Pattern == "" || cr.Severity == "" {
			return fmt.Errorf("自定义规则缺少必填字段 (id/pattern/severity)")
		}
		rule, err := NewRule(cr)
		if err != nil {
			return fmt.Errorf("自定义规则 %s: %w", cr.ID, err)
		}
		if idx, ok := existing[cr.ID]; ok {
			e.rules[idx] = rule
		} else {
			e.rules = append(e.rules, rule)
		}
	}
	return nil
}
