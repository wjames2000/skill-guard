package scanner

import (
	"os"
	"path/filepath"
	"testing"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func TestScan_FindsRisk(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "danger.py"), []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0644)
	os.WriteFile(filepath.Join(dir, "safe.py"), []byte(`print("hello")`), 0644)

	cfg := &pkgtypes.Config{
		Paths:   []string{dir},
		MaxSize: 10 * 1024 * 1024,
	}
	report, err := Scan(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if report.TotalIssues == 0 {
		t.Error("应检测到风险")
	}
	if report.TotalFiles != 2 {
		t.Errorf("期望 2 个文件, 得到 %d", report.TotalFiles)
	}
}

func TestScan_CleanDirectory(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "safe.py"), []byte(`print("hello")`), 0644)

	cfg := &pkgtypes.Config{
		Paths:   []string{dir},
		MaxSize: 10 * 1024 * 1024,
	}
	report, err := Scan(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if report.TotalIssues != 0 {
		t.Errorf("安全目录期望 0 风险, 得到 %d", report.TotalIssues)
	}
}

func TestScan_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	cfg := &pkgtypes.Config{
		Paths:   []string{dir},
		MaxSize: 10 * 1024 * 1024,
	}
	report, err := Scan(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if report.TotalFiles != 0 {
		t.Errorf("空目录期望 TotalFiles=0, 得到 %d", report.TotalFiles)
	}
	if report.Duration == "" {
		t.Error("Duration 不应为空")
	}
}

func TestScan_JsonOutput(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "danger.py"), []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0644)

	cfg := &pkgtypes.Config{
		Paths:      []string{dir},
		MaxSize:    10 * 1024 * 1024,
		JSONOutput: true,
	}
	report, err := Scan(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if report.TotalIssues == 0 {
		t.Error("JSON 模式下也应检测到风险")
	}
}

func TestScan_WithSeverityFilter(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "danger.py"), []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0644)

	cfg := &pkgtypes.Config{
		Paths:    []string{dir},
		MaxSize:  10 * 1024 * 1024,
		Severity: "critical",
	}
	report, err := Scan(cfg)
	if err != nil {
		t.Fatal(err)
	}
	_ = report
}

func TestScan_WithIgnore(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	os.MkdirAll(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "secret.py"), []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0644)
	os.WriteFile(filepath.Join(dir, "real.py"), []byte(`print("safe")`), 0644)

	cfg := &pkgtypes.Config{
		Paths:   []string{dir},
		MaxSize: 10 * 1024 * 1024,
	}
	report, err := Scan(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if report.TotalIssues != 0 {
		t.Error(".git 目录中的风险应被忽略")
	}
	if report.TotalFiles != 1 {
		t.Errorf("期望 1 个文件（real.py）, 得到 %d", report.TotalFiles)
	}
}

func TestScan_InvalidPath(t *testing.T) {
	cfg := &pkgtypes.Config{
		Paths:   []string{"/nonexistent/path"},
		MaxSize: 10 * 1024 * 1024,
	}
	report, err := Scan(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if report.TotalFiles != 0 {
		t.Errorf("无效路径期望 TotalFiles=0, 得到 %d", report.TotalFiles)
	}
}
