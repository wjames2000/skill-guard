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

	subcommand := "rules"
	if len(args) > 1 {
		subcommand = args[1]
	}

	switch subcommand {
	case "rules":
		fmt.Println("正在更新规则库...")
		rulesDir := getRulesDir()
		if err := updater.UpdateRules(rulesDir); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: 规则更新失败: %v\n", err)
			os.Exit(2)
		}
		fmt.Println("✅ 规则库已更新到最新版本")
	case "check":
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
	default:
		fmt.Fprintf(os.Stderr, "用法: skill-guard update [rules|check]\n")
		os.Exit(2)
	}
	os.Exit(0)
}

func getRulesDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".skill-guard/rules"
	}
	return filepath.Join(home, ".config", "skill-guard", "rules")
}
