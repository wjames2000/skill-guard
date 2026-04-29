package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wjames2000/skill-guard/internal/updater"
)

func handleUpdate(args []string) {
	if len(args) == 0 || args[0] != "update" {
		return
	}

	subcommand := "check"
	if len(args) > 1 {
		subcommand = args[1]
	}

	switch subcommand {
	case "check":
		handleUpdateCheck()
	case "rules":
		handleUpdateRules()
	case "list":
		handleRuleList()
	case "info":
		handleRuleInfo()
	case "install":
		handleRuleInstall(args)
	case "index":
		handleUpdateIndex()
	case "versions":
		handleRuleVersions()
	default:
		fmt.Fprintf(os.Stderr, "用法:\n")
		fmt.Fprintf(os.Stderr, "  skill-guard update check     — 检查工具版本更新\n")
		fmt.Fprintf(os.Stderr, "  skill-guard update rules     — 更新官方规则库\n")
		fmt.Fprintf(os.Stderr, "  skill-guard update list      — 列出已安装规则\n")
		fmt.Fprintf(os.Stderr, "  skill-guard update info      — 查看可用规则源\n")
		fmt.Fprintf(os.Stderr, "  skill-guard update install <url> — 安装远程规则\n")
		fmt.Fprintf(os.Stderr, "  skill-guard update index     — 刷新规则市场索引\n")
		fmt.Fprintf(os.Stderr, "  skill-guard update versions  — 查看规则版本历史\n")
		os.Exit(2)
	}
	os.Exit(0)
}

func handleUpdateCheck() {
	fmt.Println("正在检查新版本...")
	release, err := updater.CheckForUpdates()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: 检查更新失败: %v\n", err)
		os.Exit(2)
	}
	fmt.Printf("最新版本: %s\n", release.TagName)
	fmt.Printf("发布说明: %s\n", release.HTMLURL)
	if updater.IsNewerVersion(release.TagName, updater.AppVersion) {
		fmt.Println("📦 有新版本可用！请前往 GitHub Releases 页面下载。")
	} else {
		fmt.Println("✅ 当前已是最新版本")
	}
}

func handleUpdateRules() {
	fmt.Println("正在更新官方规则库...")
	rulesDir := getRulesDir()
	if err := updater.UpdateRulesFromSource(rulesDir, updater.DefaultRulesURL()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: 规则更新失败: %v\n", err)
		os.Exit(2)
	}
	fmt.Println("✅ 规则库已更新到最新版本")
}

func handleRuleList() {
	rulesDir := getRulesDir()
	files, err := updater.ListInstalledRules(rulesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}
	if len(files) == 0 {
		fmt.Println("未安装任何规则文件。运行 'skill-guard update rules' 下载官方规则。")
		return
	}
	fmt.Println("已安装的规则文件:")
	for _, f := range files {
		fmt.Printf("  📄 %s\n", f)
	}
}

func handleRuleInfo() {
	fmt.Println("正在查询规则市场...")
	index, err := updater.FetchRulesIndex()
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 无法获取远程索引: %v\n", err)
		fmt.Println("使用默认索引...")
		index = updater.DefaultRulesIndex()
	}
	fmt.Printf("规则市场索引 (v%s, 更新于 %s):\n", index.Version, index.Updated)
	for _, s := range index.Sources {
		fmt.Printf("  📦 %-12s %-20s (%d 条规则)\n", s.ID, s.Name, s.Rules)
		fmt.Printf("     %s\n", s.Description)
	}
}

func handleRuleInstall(args []string) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "用法: skill-guard update install <rule-url>\n")
		os.Exit(2)
	}
	url := args[2]
	fmt.Printf("正在安装规则: %s\n", url)
	rulesDir := getRulesDir()
	if err := updater.UpdateRulesFromSource(rulesDir, url); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: 安装失败: %v\n", err)
		os.Exit(2)
	}
	fmt.Println("✅ 规则安装完成")
}

func handleRuleVersions() {
	rulesDir := getRulesDir()
	versions, err := updater.ListVersions(rulesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}
	if len(versions) == 0 {
		fmt.Println("暂无版本记录。运行 'skill-guard update rules' 后自动记录。")
		return
	}
	fmt.Println("规则版本历史:")
	for _, v := range versions {
		fmt.Printf("  📄 %-30s v%-12s 安装于 %s\n", v.File, v.Version, v.Installed[:10])
	}
}

func handleUpdateIndex() {
	fmt.Println("正在刷新规则市场索引...")
	index, err := updater.FetchRulesIndex()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}
	fmt.Printf("规则市场索引已刷新 (v%s, %d 个规则源)\n", index.Version, len(index.Sources))
}

func getRulesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".skill-guard/rules"
	}
	return filepath.Join(home, ".config", "skill-guard", "rules")
}
