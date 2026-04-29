package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration_DetectsRisk(t *testing.T) {
	build := exec.Command("go", "build", "-o", "skill-guard.test", ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("构建失败: %v\n%s", err, out)
	}
	defer os.Remove("skill-guard.test")

	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "danger.py"), []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0644)
	os.WriteFile(filepath.Join(dir, "safe.py"), []byte(`print("hello")`), 0644)

	cmd := exec.Command("./skill-guard.test", "--json", dir)
	output, err := cmd.Output()
	var exitErr *exec.ExitError
	if err != nil {
		exitErr, _ = err.(*exec.ExitError)
	}
	if exitErr != nil && exitErr.ExitCode() != 1 {
		t.Fatalf("退出码应为 1, 得到 %d: %s", exitErr.ExitCode(), string(output))
	}
	if !strings.Contains(string(output), "SKL-001") {
		t.Error("应检测到 SKL-001")
	}
}

func TestIntegration_CleanFile(t *testing.T) {
	build := exec.Command("go", "build", "-o", "skill-guard.test", ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("构建失败: %v\n%s", err, out)
	}
	defer os.Remove("skill-guard.test")

	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "safe.py"), []byte(`print("hello")`), 0644)

	cmd := exec.Command("./skill-guard.test", "--json", dir)
	output, err := cmd.Output()
	if err != nil {
		t.Fatal("扫描失败:", err)
	}
	if !strings.Contains(string(output), `"total_issues": 0`) {
		t.Error("安全文件应无风险")
	}
}
