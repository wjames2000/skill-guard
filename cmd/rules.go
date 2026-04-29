package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
	"github.com/wjames2000/skill-guard/internal/engine"
)

func handleRules(args []string) {
	if len(args) == 0 || args[0] != "rules" {
		return
	}

	subcommand := "help"
	if len(args) > 1 {
		subcommand = args[1]
	}

	switch subcommand {
	case "export":
		handleRulesExport(args)
	case "new":
		handleRulesNew(args)
	case "test":
		handleRulesTest(args)
	case "help", "--help", "-h":
		printRulesHelp()
	default:
		printRulesHelp()
		os.Exit(2)
	}
	os.Exit(0)
}

func printRulesHelp() {
	fmt.Println("用法:")
	fmt.Println("  skill-guard rules export [format]  — 导出内置规则 (yaml/json)")
	fmt.Println("  skill-guard rules new              — 创建新规则模板")
	fmt.Println("  skill-guard rules test <file> [input]  — 测试规则匹配")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  skill-guard rules export           — 导出为 YAML 到 stdout")
	fmt.Println("  skill-guard rules export json      — 导出为 JSON 到 stdout")
	fmt.Println("  skill-guard rules export > my-rules.yaml — 保存到文件")
	fmt.Println("  skill-guard rules new              — 生成规则模板")
	fmt.Println("  skill-guard rules test my-rule.yaml \"dangerous_code()\"  — 测试规则")
}

func handleRulesExport(args []string) {
	format := "yaml"
	if len(args) > 2 {
		format = args[2]
	}

	rules := engine.BuiltinRules()

	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(rules); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: 序列化失败: %v\n", err)
			os.Exit(2)
		}
	case "yaml", "yml":
		data, err := yaml.Marshal(struct {
			Rules []*pkgtypes.Rule `yaml:"rules"`
		}{Rules: rules})
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: 序列化失败: %v\n", err)
			os.Exit(2)
		}
		fmt.Println(string(data))
	default:
		fmt.Fprintf(os.Stderr, "不支持的格式: %s (支持: yaml, json)\n", format)
		os.Exit(2)
	}
}

func handleRulesNew(args []string) {
	template := pkgtypes.Rule{
		ID:          "USR-001",
		Name:        "新规则名称",
		Severity:    "Medium",
		Description: "规则的详细说明",
		Pattern:     "regex_pattern_here",
		Keywords:    []string{"keyword1", "keyword2"},
		FileTypes:   []string{".py", ".sh"},
		Ref:         "参考链接（可选）",
	}

	data, err := yaml.Marshal(struct {
		Rules []pkgtypes.Rule `yaml:"rules"`
	}{Rules: []pkgtypes.Rule{template}})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: 生成模板失败: %v\n", err)
		os.Exit(2)
	}

	fmt.Println("# 规则模板 - 编辑后保存为 .yaml 文件")
	fmt.Println("# 加载方式: skill-guard --rules my-rule.yaml")
	fmt.Println("# 参考: https://github.com/wjames2000/skill-guard/wiki")
	fmt.Println(string(data))

	filename := fmt.Sprintf("custom-rule-%s.yaml", time.Now().Format("20060102-150405"))
	if err := os.WriteFile(filename, data, 0644); err == nil {
		fmt.Fprintf(os.Stderr, "✅ 模板已保存到 %s\n", filename)
	}
}

func handleRulesTest(args []string) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "用法: skill-guard rules test <rule-file.yaml> [sample-input]\n")
		fmt.Fprintf(os.Stderr, "示例: skill-guard rules test my-rule.yaml \"dangerous_code()\"\n")
		os.Exit(2)
	}

	ruleFile := args[2]
	sampleInput := ""
	if len(args) > 3 {
		sampleInput = args[3]
	}

	rules, err := engine.LoadRulesFile(ruleFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: 加载规则失败: %v\n", err)
		os.Exit(2)
	}

	for _, r := range rules {
		rule, err := engine.NewRule(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: 编译规则 %s 失败: %v\n", r.ID, err)
			continue
		}

		fmt.Printf("规则: %s [%s] %s\n", r.ID, r.Severity, r.Name)
		fmt.Printf("  模式: %s\n", r.Pattern)

		if sampleInput != "" {
			result := rule.MatchLine(sampleInput, "<input>", 1)
			if result != nil {
				fmt.Printf("  ✅ 匹配成功 (行 %d)\n", result.LineNumber)
			} else {
				fmt.Printf("  ❌ 未匹配\n")
			}
		}

		if len(r.Keywords) > 0 {
			hasKW := rule.HasKeywords([]string{sampleInput})
			fmt.Printf("  关键词过滤: %v\n", hasKW)
		}
	}
}
