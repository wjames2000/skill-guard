package output

import (
	"fmt"
	"io"
	"os"
	"runtime"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
)

var ColorEnabled = true

func init() {
	if runtime.GOOS == "windows" {
		ColorEnabled = false
	} else {
		fi, err := os.Stdout.Stat()
		if err != nil || (fi.Mode()&os.ModeCharDevice) == 0 {
			ColorEnabled = false
		}
	}
}

func c(code string) string {
	if !ColorEnabled {
		return ""
	}
	return code
}

type TerminalRenderer struct{}

func (t *TerminalRenderer) Render(w io.Writer, report *pkgtypes.ScanReport) error {
	fmt.Fprintf(w, "╔════════════════════════════════════════════════════════╗\n")
	fmt.Fprintf(w, "║  %sskill-guard%s  扫描报告                              ║\n", c(colorBold), c(colorReset))
	fmt.Fprintf(w, "║  扫描时间: %-42s ║\n", report.ScanTime)
	fmt.Fprintf(w, "║  扫描文件数: %-41d ║\n", report.TotalFiles)
	fmt.Fprintf(w, "║  发现风险: %-43d ║\n", report.TotalIssues)
	fmt.Fprintf(w, "╚════════════════════════════════════════════════════════╝\n\n")

	if len(report.Results) == 0 {
		fmt.Fprintf(w, "%s✓ 未发现安全风险%s\n", c(colorGreen), c(colorReset))
		return nil
	}

	for _, result := range report.Results {
		col := severityColor(result.Severity)
		fmt.Fprintf(w, "%s%s [%s] %s%s\n", c(col), result.Severity, result.RuleID, result.RuleName, c(colorReset))
		fmt.Fprintf(w, "  → %s:%d\n", result.FilePath, result.LineNumber)
		if result.LineContent != "" {
			fmt.Fprintf(w, "    %s\n", result.LineContent)
		}
	}

	fmt.Fprintf(w, "\n%s━━━━━━━━━━ 汇总 ━━━━━━━━━━%s\n", c(colorBold), c(colorReset))
	fmt.Fprintf(w, "  Critical: %d | High: %d | Medium: %d | Low: %d\n",
		report.Summary.Critical, report.Summary.High,
		report.Summary.Medium, report.Summary.Low)
	return nil
}

func severityColor(severity string) string {
	switch severity {
	case "Critical":
		return colorRed + colorBold
	case "High":
		return colorRed
	case "Medium":
		return colorYellow
	case "Low":
		return colorBlue
	default:
		return colorCyan
	}
}

