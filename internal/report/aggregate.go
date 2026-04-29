package report

import (
	"sort"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

var severityOrder = map[string]int{
	"Critical": 0,
	"High":     1,
	"Medium":   2,
	"Low":      3,
}

func Aggregate(results []*pkgtypes.MatchResult) []*pkgtypes.MatchResult {
	seen := make(map[string]bool)
	var unique []*pkgtypes.MatchResult
	for _, r := range results {
		key := r.RuleID + ":" + r.FilePath
		if seen[key] {
			continue
		}
		seen[key] = true
		unique = append(unique, r)
	}
	sort.Slice(unique, func(i, j int) bool {
		return severityOrder[unique[i].Severity] < severityOrder[unique[j].Severity]
	})
	return unique
}
