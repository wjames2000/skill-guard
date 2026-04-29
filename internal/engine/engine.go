package engine

import (
	"fmt"

	"github.com/wjames2000/skill-guard/internal/file"
	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

type Engine struct {
	rules []*Rule
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

func (e *Engine) Match(target *pkgtypes.FileTarget) []*pkgtypes.MatchResult {
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
