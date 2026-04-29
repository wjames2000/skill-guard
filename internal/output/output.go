package output

import (
	"fmt"
	"io"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

type Renderer interface {
	Render(w io.Writer, report *pkgtypes.ScanReport) error
}

func Render(w io.Writer, report *pkgtypes.ScanReport, jsonOutput, quiet, summary, sarif bool) error {
	var renderer Renderer
	switch {
	case jsonOutput:
		renderer = &JSONRenderer{}
	case quiet:
		renderer = &QuietRenderer{}
	case summary:
		renderer = &SummaryRenderer{}
	case sarif:
		renderer = &SARIFRenderer{}
	default:
		renderer = &TerminalRenderer{}
	}
	return renderer.Render(w, report)
}

type SummaryRenderer struct{}

func (s *SummaryRenderer) Render(w io.Writer, report *pkgtypes.ScanReport) error {
	fmt.Fprintf(w, "扫描文件: %d\n", report.TotalFiles)
	fmt.Fprintf(w, "发现风险: %d\n", report.TotalIssues)
	fmt.Fprintf(w, "  Critical: %d\n", report.Summary.Critical)
	fmt.Fprintf(w, "  High:     %d\n", report.Summary.High)
	fmt.Fprintf(w, "  Medium:   %d\n", report.Summary.Medium)
	fmt.Fprintf(w, "  Low:      %d\n", report.Summary.Low)
	if report.Duration != "" {
		fmt.Fprintf(w, "扫描耗时: %s\n", report.Duration)
	}
	return nil
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
