package engine

import (
	"fmt"
	"regexp"
	"strings"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

type Rule struct {
	*pkgtypes.Rule
	compiled *regexp.Regexp
}

func NewRule(r *pkgtypes.Rule) (*Rule, error) {
	compiled, err := regexp.Compile(r.Pattern)
	if err != nil {
		return nil, fmt.Errorf("规则 %s: 无效正则: %w", r.ID, err)
	}
	return &Rule{Rule: r, compiled: compiled}, nil
}

func (r *Rule) MatchesFileType(ext string) bool {
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
