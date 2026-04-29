package report

import (
	"time"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func Build(results []*pkgtypes.MatchResult, totalFiles int, start time.Time) *pkgtypes.ScanReport {
	results = Aggregate(results)
	return &pkgtypes.ScanReport{
		ScanTime:    start.Format("2006-01-02T15:04:05-07:00"),
		Duration:    time.Since(start).String(),
		TotalFiles:  totalFiles,
		TotalIssues: len(results),
		Results:     results,
		Summary:     CalculateSummary(results),
	}
}
