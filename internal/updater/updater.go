package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultRulesURL = "https://github.com/wjames2000/skill-guard/releases/latest/download/builtin_rules.yaml"
	updateCheckURL  = "https://api.github.com/repos/wjames2000/skill-guard/releases/latest"
	AppVersion      = "v0.1.0"
)

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Body    string `json:"body"`
}

func CheckForUpdates() (*ReleaseInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(updateCheckURL)
	if err != nil {
		return nil, fmt.Errorf("检查更新失败: %w", err)
	}
	defer resp.Body.Close()

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("解析更新信息失败: %w", err)
	}
	return &release, nil
}

func UpdateRules(rulesDir string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(defaultRulesURL)
	if err != nil {
		return fmt.Errorf("下载规则失败: %w", err)
	}
	defer resp.Body.Close()

	os.MkdirAll(rulesDir, 0755)
	dst := filepath.Join(rulesDir, "builtin_rules.yaml")
	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建规则文件失败: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("写入规则文件失败: %w", err)
	}
	return nil
}

func IsNewerVersion(remote, local string) bool {
	return remote > local
}
