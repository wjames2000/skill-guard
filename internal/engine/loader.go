package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	pkgtypes "github.com/hpds.cc/skill-guard/pkg/types"
)

func LoadRulesFile(path string) ([]*pkgtypes.Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取规则文件失败: %w", err)
	}
	ext := strings.ToLower(filepath.Ext(path))
	var rules []*pkgtypes.Rule
	switch ext {
	case ".yaml", ".yml":
		var list pkgtypes.RuleList
		if err := yaml.Unmarshal(data, &list); err != nil {
			return nil, fmt.Errorf("YAML 解析失败: %w", err)
		}
		for i := range list.Rules {
			rules = append(rules, &list.Rules[i])
		}
	case ".json":
		if err := json.Unmarshal(data, &rules); err != nil {
			return nil, fmt.Errorf("JSON 解析失败: %w", err)
		}
	default:
		return nil, fmt.Errorf("不支持的规则文件格式: %s", ext)
	}
	if len(rules) == 0 {
		return nil, fmt.Errorf("规则文件为空")
	}
	return rules, nil
}
