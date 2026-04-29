package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func TestTerminalRenderer_NoIssues(t *testing.T) {
	report := &pkgtypes.ScanReport{TotalFiles: 10, TotalIssues: 0, Results: []*pkgtypes.MatchResult{}, Summary: &pkgtypes.Summary{}}
	var buf bytes.Buffer
	(&TerminalRenderer{}).Render(&buf, report)
	if !strings.Contains(buf.String(), "未发现安全风险") {
		t.Error("无风险时应显示提示信息")
	}
}

func TestTerminalRenderer_WithIssues(t *testing.T) {
	report := &pkgtypes.ScanReport{
		TotalFiles: 1, TotalIssues: 1,
		Results: []*pkgtypes.MatchResult{{RuleID: "SKL-001", Severity: "Critical", FilePath: "test.py", LineNumber: 15, LineContent: "bad code"}},
		Summary: &pkgtypes.Summary{Critical: 1},
	}
	var buf bytes.Buffer
	(&TerminalRenderer{}).Render(&buf, report)
	output := buf.String()
	if !strings.Contains(output, "SKL-001") {
		t.Error("应包含规则 ID")
	}
	if !strings.Contains(output, "test.py:15") {
		t.Error("应包含文件路径和行号")
	}
}

func TestJSONRenderer(t *testing.T) {
	report := &pkgtypes.ScanReport{TotalFiles: 5, TotalIssues: 1, Results: []*pkgtypes.MatchResult{{RuleID: "SKL-001", Severity: "Critical", FilePath: "test.py", LineNumber: 10}}, Summary: &pkgtypes.Summary{Critical: 1}}
	var buf bytes.Buffer
	(&JSONRenderer{}).Render(&buf, report)
	var decoded pkgtypes.ScanReport
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatal("JSON 不合法:", err)
	}
	if decoded.TotalIssues != 1 {
		t.Errorf("TotalIssues 期望 1, 得到 %d", decoded.TotalIssues)
	}
}

func TestQuietRenderer(t *testing.T) {
	report := &pkgtypes.ScanReport{
		Results: []*pkgtypes.MatchResult{
			{FilePath: "a.py"}, {FilePath: "a.py"}, {FilePath: "b.py"},
		},
	}
	var buf bytes.Buffer
	(&QuietRenderer{}).Render(&buf, report)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("期望 2 行去重输出, 得到 %d", len(lines))
	}
}

func TestSeverityFilter(t *testing.T) {
	results := []*pkgtypes.MatchResult{
		{Severity: "Critical"}, {Severity: "High"}, {Severity: "Medium"}, {Severity: "Low"},
	}
	filtered := SeverityFilter(results, "high")
	if len(filtered) != 2 {
		t.Errorf("high 过滤期望 2 条(Critical+High), 得到 %d", len(filtered))
	}

	filtered2 := SeverityFilter(results, "")
	if len(filtered2) != 4 {
		t.Errorf("空过滤期望 4 条, 得到 %d", len(filtered2))
	}
}

func TestSeverityFilter_Empty(t *testing.T) {
	var empty []*pkgtypes.MatchResult
	filtered := SeverityFilter(empty, "high")
	if len(filtered) != 0 {
		t.Errorf("空结果应返回空, 得到 %d", len(filtered))
	}
	filtered2 := SeverityFilter(empty, "")
	if len(filtered2) != 0 {
		t.Errorf("空结果+空过滤应返回空, 得到 %d", len(filtered2))
	}
}

func TestSummaryRenderer(t *testing.T) {
	report := &pkgtypes.ScanReport{
		TotalFiles: 100, TotalIssues: 5,
		Summary: &pkgtypes.Summary{Critical: 2, High: 2, Medium: 1, Low: 0},
		Duration: "1.5s",
	}
	var buf bytes.Buffer
	r := &SummaryRenderer{}
	if err := r.Render(&buf, report); err != nil {
		t.Fatal(err)
	}
	output := buf.String()
	if !strings.Contains(output, "100") {
		t.Error("应包含总文件数")
	}
	if !strings.Contains(output, "5") {
		t.Error("应包含总风险数")
	}
	if !strings.Contains(output, "1.5s") {
		t.Error("应包含耗时")
	}
}

func TestSeverityFilter_InvalidLevel(t *testing.T) {
	results := []*pkgtypes.MatchResult{{Severity: "Critical"}, {Severity: "High"}}
	filtered := SeverityFilter(results, "invalid")
	if len(filtered) != 2 {
		t.Errorf("无效级别应返回全部，得到 %d", len(filtered))
	}
}

func TestQuietRenderer_Empty(t *testing.T) {
	report := &pkgtypes.ScanReport{Results: []*pkgtypes.MatchResult{}}
	var buf bytes.Buffer
	(&QuietRenderer{}).Render(&buf, report)
	if buf.Len() != 0 {
		t.Error("空结果 quiet 模式应无输出")
	}
}
