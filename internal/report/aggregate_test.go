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

func TestAggregate_EmptyResults(t *testing.T) {
	agg := Aggregate(nil)
	if len(agg) != 0 {
		t.Errorf("空结果应返回空, 得到 %d", len(agg))
	}
	agg2 := Aggregate([]*pkgtypes.MatchResult{})
	if len(agg2) != 0 {
		t.Errorf("空切片应返回空, 得到 %d", len(agg2))
	}
}

func TestAggregate_SameSeverity(t *testing.T) {
	results := []*pkgtypes.MatchResult{
		{RuleID: "SKL-001", FilePath: "a.py", Severity: "High"},
		{RuleID: "SKL-002", FilePath: "b.py", Severity: "High"},
		{RuleID: "SKL-003", FilePath: "c.py", Severity: "High"},
	}
	agg := Aggregate(results)
	if len(agg) != 3 {
		t.Errorf("同级别 3 条去重后期望 3, 得到 %d", len(agg))
	}
}
