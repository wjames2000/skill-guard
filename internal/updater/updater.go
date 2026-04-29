package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultRulesURL = "https://github.com/wjames2000/skill-guard/releases/latest/download/builtin_rules.yaml"
	updateCheckURL  = "https://api.github.com/repos/wjames2000/skill-guard/releases/latest"
	rulesIndexURL   = "https://github.com/wjames2000/skill-guard/releases/latest/download/community-index.yaml"
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

type RuleSource struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	URL         string `yaml:"url" json:"url"`
	Rules       int    `yaml:"rules" json:"rules"`
}

type RulesIndex struct {
	Version string       `yaml:"version" json:"version"`
	Updated string       `yaml:"updated" json:"updated"`
	Sources []RuleSource `yaml:"sources" json:"sources"`
}

func FetchRulesIndex() (*RulesIndex, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(rulesIndexURL)
	if err != nil {
		return nil, fmt.Errorf("获取规则索引失败: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取规则索引失败: %w", err)
	}

	var index RulesIndex
	if err := yaml.Unmarshal(data, &index); err != nil {
		return DefaultRulesIndex(), nil
	}
	return &index, nil
}

func DefaultRulesIndex() *RulesIndex {
	return &RulesIndex{
		Version: "1",
		Updated: time.Now().Format("2006-01-02"),
		Sources: []RuleSource{
			{ID: "official", Name: "官方规则库", Description: "skill-guard 官方维护的安全检测规则", URL: defaultRulesURL, Rules: 32},
		},
	}
}

func DefaultRulesURL() string {
	return defaultRulesURL
}

func UpdateRules(rulesDir string) error {
	return UpdateRulesFromSource(rulesDir, defaultRulesURL)
}

func UpdateRulesFromSource(rulesDir, sourceURL string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(sourceURL)
	if err != nil {
		return fmt.Errorf("下载规则失败: %w", err)
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(rulesDir, 0750); err != nil {
		return fmt.Errorf("创建规则目录失败: %w", err)
	}

	filename := filepath.Base(sourceURL)
	if filename == "" || !strings.HasSuffix(filename, ".yaml") {
		filename = "rules.yaml"
	}
	dst := filepath.Join(rulesDir, filename)
	f, err := os.Create(dst) // #nosec G304
	if err != nil {
		return fmt.Errorf("创建规则文件失败: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("写入规则文件失败: %w", err)
	}
	return nil
}

func ListInstalledRules(rulesDir string) ([]string, error) {
	entries, err := os.ReadDir(rulesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".yaml") {
			files = append(files, e.Name())
		}
	}
	return files, nil
}

func IsNewerVersion(remote, local string) bool {
	return remote > local
}
