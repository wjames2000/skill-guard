package output

import (
	"io"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

type Renderer interface {
	Render(w io.Writer, report *pkgtypes.ScanReport) error
}

func Render(w io.Writer, report *pkgtypes.ScanReport, jsonOutput, quiet bool) error {
	var renderer Renderer
	switch {
	case jsonOutput:
		renderer = &JSONRenderer{}
	case quiet:
		renderer = &QuietRenderer{}
	default:
		renderer = &TerminalRenderer{}
	}
	return renderer.Render(w, report)
}

func SeverityFilter(results []*pkgtypes.MatchResult, severity string) []*pkgtypes.MatchResult {
	if severity == "" {
		return results
	}
	weight := map[string]int{"Critical": 4, "High": 3, "Medium": 2, "Low": 1}
	threshold, ok := map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1}[severity]
	if !ok {
		return results
	}
	var filtered []*pkgtypes.MatchResult
	for _, r := range results {
		if weight[r.Severity] >= threshold {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
