package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func handleInit(args []string) {
	if len(args) == 0 || args[0] != "init" {
		return
	}

	fmt.Println("🔧 初始化 skill-guard 项目配置...")

	if _, err := os.Stat(".skillguard.yaml"); os.IsNotExist(err) {
		defaultConfig := `# skill-guard 配置文件
severity: "medium"
ext-include:
  - ".py"
  - ".sh"
  - ".js"
  - ".yaml"
  - ".json"
`
		if err := os.WriteFile(".skillguard.yaml", []byte(defaultConfig), 0644); err != nil { // #nosec G306
			fmt.Fprintf(os.Stderr, "ERROR: 创建配置文件失败: %v\n", err)
			os.Exit(2)
		}
		fmt.Println("  ✅ 已创建 .skillguard.yaml")
	} else {
		fmt.Println("  ⏭️  .skillguard.yaml 已存在，跳过")
	}

	gitDir := ".git"
	if info, err := os.Stat(gitDir); err != nil || !info.IsDir() {
		fmt.Println("  ⏭️  未检测到 Git 仓库，跳过 pre-commit hook 安装")
		fmt.Println("  提示: 在 Git 仓库中运行 'skill-guard init' 可自动安装 pre-commit hook")
		os.Exit(0)
	}

	hooksDir := filepath.Join(gitDir, "hooks")
	hookPath := filepath.Join(hooksDir, "pre-commit")

	if _, err := os.Stat(hookPath); err == nil {
		fmt.Println("  ⏭️  pre-commit hook 已存在，跳过")
		fmt.Println("  如需覆盖，请手动删除 .git/hooks/pre-commit 后重试")
		os.Exit(0)
	}

	hookContent := `#!/bin/sh
# skill-guard pre-commit hook — 自动安装 by skill-guard init
echo "🔍 skill-guard: 扫描暂存文件中的安全风险..."
changed_files=$(git diff --cached --name-only --diff-filter=ACM | grep -E '\.(py|sh|js|ts|yaml|yml|json|toml|md|go)$' || true)

if [ -z "$changed_files" ]; then
    echo "✅ skill-guard: 无可扫描文件，跳过"
    exit 0
fi

skill-guard $changed_files --severity high --quiet
result=$?

if [ $result -eq 1 ]; then
    echo "❌ skill-guard: 发现安全风险，请修复后重新提交"
    echo "   使用 --no-verify 可跳过检查（不推荐）"
elif [ $result -eq 0 ]; then
    echo "✅ skill-guard: 安全扫描通过"
fi

exit $result
`
	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil { // #nosec G306
		fmt.Fprintf(os.Stderr, "ERROR: 安装 pre-commit hook 失败: %v\n", err)
		os.Exit(2)
	}
	fmt.Println("  ✅ 已安装 pre-commit hook (.git/hooks/pre-commit)")
	fmt.Println("📦 初始化完成！后续提交将自动扫描安全风险。")
	os.Exit(0)
}
