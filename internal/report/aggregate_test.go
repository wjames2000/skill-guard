package report

import (
	"testing"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func TestAggregate_Deduplicate(t *testing.T) {
	results := []*pkgtypes.MatchResult{
		{RuleID: "SKL-001", FilePath: "a.py", Severity: "High"},
		{RuleID: "SKL-001", FilePath: "a.py", Severity: "High"},
		{RuleID: "SKL-002", FilePath: "a.py", Severity: "Low"},
	}
	agg := Aggregate(results)
	if len(agg) != 2 {
		t.Errorf("期望 2 条去重后结果, 得到 %d", len(agg))
	}
}

func TestAggregate_Order(t *testing.T) {
	results := []*pkgtypes.MatchResult{
		{RuleID: "SKL-003", FilePath: "a.py", Severity: "Low"},
		{RuleID: "SKL-001", FilePath: "a.py", Severity: "Critical"},
		{RuleID: "SKL-002", FilePath: "a.py", Severity: "Medium"},
	}
	agg := Aggregate(results)
	expected := []string{"Critical", "Medium", "Low"}
	for i, r := range agg {
		if r.Severity != expected[i] {
			t.Errorf("位置 %d: 期望 %s, 得到 %s", i, expected[i], r.Severity)
		}
	}
}
