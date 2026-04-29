package report

import (
	"testing"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func TestCalculateSummary(t *testing.T) {
	results := []*pkgtypes.MatchResult{
		{Severity: "Critical"},
		{Severity: "High"},
		{Severity: "High"},
		{Severity: "Medium"},
	}
	s := CalculateSummary(results)
	if s.Critical != 1 {
		t.Errorf("Critical 期望 1, 得到 %d", s.Critical)
	}
	if s.High != 2 {
		t.Errorf("High 期望 2, 得到 %d", s.High)
	}
	if s.Medium != 1 {
		t.Errorf("Medium 期望 1, 得到 %d", s.Medium)
	}
	if s.Low != 0 {
		t.Errorf("Low 期望 0, 得到 %d", s.Low)
	}
}
