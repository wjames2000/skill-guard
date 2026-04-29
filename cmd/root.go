package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wjames2000/skill-guard/internal/config"
	"github.com/wjames2000/skill-guard/internal/output"
	"github.com/wjames2000/skill-guard/internal/report"
	"github.com/wjames2000/skill-guard/internal/scanner"
	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func Execute() {
	cfg := parseFlags()

	configPath := cfg.ConfigFile
	if configPath == "" {
		configPath = ".skillguard.yaml"
	}
	fileCfg, err := config.LoadFile(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}
	cfg = config.MergeWithCLI(cfg, fileCfg)

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}
	if err := runScan(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}
}

func parseFlags() *pkgtypes.Config {
	cfg := pkgtypes.DefaultConfig()
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json", "-j":
			cfg.JSONOutput = true
		case "--quiet", "-q":
			cfg.Quiet = true
		case "--verbose", "-v":
			cfg.Verbose = true
		case "--config", "-c":
			i++
			if i < len(args) {
				cfg.ConfigFile = args[i]
			}
		case "--rules", "-r":
			i++
			if i < len(args) {
				cfg.RulesFile = args[i]
			}
		case "--severity", "-s":
			i++
			if i < len(args) {
				cfg.Severity = args[i]
			}
		case "--ignore", "-i":
			i++
			if i < len(args) {
				cfg.Ignore = append(cfg.Ignore, args[i])
			}
		case "--ext-include":
			i++
			if i < len(args) {
				cfg.ExtInclude = splitExt(args[i])
			}
		case "--ext-exclude":
			i++
			if i < len(args) {
				cfg.ExtExclude = splitExt(args[i])
			}
		case "--max-size":
			i++
			if i < len(args) {
				if v, err := parseMaxSizeFlag(args[i]); err == nil {
					cfg.MaxSize = v
				}
			}
		case "--concurrency":
			i++
			if i < len(args) {
				if v, err := strconv.Atoi(args[i]); err == nil && v > 0 {
					cfg.Concurrency = v
				}
			}
		case "--disable-builtin":
			cfg.DisableBuiltin = true
		case "--version":
			printVersion()
			os.Exit(0)
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		default:
			if !strings.HasPrefix(args[i], "-") {
				if i == 0 || strings.HasPrefix(args[i-1], "-") {
					cfg.Paths = []string{args[i]}
				} else {
					cfg.Paths = append(cfg.Paths, args[i])
				}
			}
		}
	}
	return cfg
}

func parseMaxSizeFlag(s string) (int64, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	multiplier := int64(1)
	switch {
	case strings.HasSuffix(s, "GB"):
		multiplier = 1024 * 1024 * 1024
		s = strings.TrimSuffix(s, "GB")
	case strings.HasSuffix(s, "MB"):
		multiplier = 1024 * 1024
		s = strings.TrimSuffix(s, "MB")
	case strings.HasSuffix(s, "KB"):
		multiplier = 1024
		s = strings.TrimSuffix(s, "KB")
	}
	v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return 0, err
	}
	return v * multiplier, nil
}

func splitExt(s string) []string {
	parts := strings.Split(s, ",")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
		if !strings.HasPrefix(parts[i], ".") {
			parts[i] = "." + parts[i]
		}
	}
	return parts
}

func runScan(cfg *pkgtypes.Config) error {
	result, err := scanner.Scan(cfg)
	if err != nil {
		return err
	}

	result.Results = output.SeverityFilter(result.Results, cfg.Severity)
	result.TotalIssues = len(result.Results)
	result.Summary = report.CalculateSummary(result.Results)

	output.Render(os.Stdout, result, cfg.JSONOutput, cfg.Quiet)

	if result.TotalIssues > 0 {
		os.Exit(1)
	}
	return nil
}

func printVersion() {
	fmt.Printf("skill-guard %s (commit: %s, built: %s)\n", Version, Commit, Date)
}

func printHelp() {
	fmt.Println(`skill-guard - 安全技能扫描工具

用法:
  skill-guard [path...] [flags]

参数:
  path  要扫描的文件或目录路径（默认: "."）

标志:
  -c, --config string        指定配置文件路径
  -r, --rules string         自定义规则文件（YAML/JSON）
  -s, --severity string      最低严重级别 (critical/high/medium/low)
  -j, --json                 JSON 格式输出
  -q, --quiet                安静模式（仅显示有问题文件）
  -v, --verbose              显示扫描进度
  -i, --ignore strings       额外忽略的路径或模式
      --ext-include strings  仅扫描指定扩展名（逗号分隔）
      --ext-exclude strings  排除指定扩展名（逗号分隔）
      --max-size string      文件大小上限（默认: 10MB）
      --concurrency int      扫描并发数（默认: CPU 核数）
      --disable-builtin      禁用内置规则
      --version              显示版本信息
  -h, --help                 显示帮助信息

示例:
  skill-guard ./skills
  skill-guard ./skills --json --severity high
  skill-guard --version`)
}
