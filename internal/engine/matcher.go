package engine

import (
	"strings"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

const maxLineContentLen = 120

func (r *Rule) HasKeywords(lines []string) bool {
	if len(r.Keywords) == 0 {
		return true
	}
	for _, line := range lines {
		for _, kw := range r.Keywords {
			if strings.Contains(line, kw) {
				return true
			}
		}
	}
	return false
}

func (r *Rule) MatchLine(line string, filePath string, lineNum int) *pkgtypes.MatchResult {
	if r.compiled == nil {
		return nil
	}
	loc := r.compiled.FindStringIndex(line)
	if loc == nil {
		return nil
	}
	content := strings.TrimSpace(line)
	if len(content) > maxLineContentLen {
		content = content[:maxLineContentLen] + "..."
	}
	return &pkgtypes.MatchResult{
		RuleID:      r.ID,
		Severity:    r.Severity,
		FilePath:    filePath,
		LineNumber:  lineNum,
		LineContent: content,
		MatchType:   "regex",
	}
}
