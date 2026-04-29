package report

import pkgtypes "github.com/wjames2000/skill-guard/pkg/types"

func CalculateSummary(results []*pkgtypes.MatchResult) *pkgtypes.Summary {
	s := &pkgtypes.Summary{}
	for _, r := range results {
		switch r.Severity {
		case "Critical":
			s.Critical++
		case "High":
			s.High++
		case "Medium":
			s.Medium++
		case "Low":
			s.Low++
		}
	}
	return s
}
