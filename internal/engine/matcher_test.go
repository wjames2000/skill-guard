package engine

import (
	"testing"

	pkgtypes "github.com/hpds.cc/skill-guard/pkg/types"
)

func TestHasKeywords(t *testing.T) {
	r, _ := NewRule(&pkgtypes.Rule{
		ID: "TEST-010", Pattern: `test`, Keywords: []string{"AKIA", "SECRET"},
	})
	if r.HasKeywords([]string{"nothing here"}) {
		t.Error("无关键字文件应返回 false")
	}
	if !r.HasKeywords([]string{"contains AKIA here"}) {
		t.Error("包含 AKIA 的文件应返回 true")
	}
	r2, _ := NewRule(&pkgtypes.Rule{ID: "TEST-011", Pattern: `test`})
	if !r2.HasKeywords([]string{"anything"}) {
		t.Error("无关键字应跳过预过滤返回 true")
	}
}

func TestMatchLine(t *testing.T) {
	r, _ := NewRule(&pkgtypes.Rule{
		ID: "SKL-001", Name: "Test AWS Key",
		Pattern: `(?i)AKIA[A-Z0-9]{16}`, Severity: "Critical",
	})
	result := r.MatchLine(`key = "AKIAIOSFODNN7EXAMPLE"`, "test.py", 10)
	if result == nil {
		t.Fatal("应匹配成功")
	}
	if result.RuleID != "SKL-001" {
		t.Errorf("RuleID 期望 SKL-001, 得到 %s", result.RuleID)
	}
	if result.LineNumber != 10 {
		t.Errorf("LineNumber 期望 10, 得到 %d", result.LineNumber)
	}
	if result.FilePath != "test.py" {
		t.Errorf("FilePath 期望 test.py, 得到 %s", result.FilePath)
	}
	if result.Severity != "Critical" {
		t.Errorf("Severity 期望 Critical, 得到 %s", result.Severity)
	}

	result2 := r.MatchLine("no key here", "test.py", 1)
	if result2 != nil {
		t.Error("不匹配的行应返回 nil")
	}
}
