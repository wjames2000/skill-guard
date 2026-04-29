package config

import (
	"strconv"
	"strings"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func MergeWithCLI(cfg *pkgtypes.Config, fileCfg *FileConfig) *pkgtypes.Config {
	cfg.Ignore = append(fileCfg.Ignore, cfg.Ignore...)
	if cfg.RulesFile == "" {
		cfg.RulesFile = fileCfg.Rules
	}
	if cfg.Severity == "" {
		cfg.Severity = fileCfg.Severity
	}
	if len(cfg.ExtInclude) == 0 {
		cfg.ExtInclude = fileCfg.ExtInclude
	}
	if len(cfg.ExtExclude) == 0 {
		cfg.ExtExclude = fileCfg.ExtExclude
	}
	if cfg.Concurrency == 0 && fileCfg.Concurrency > 0 {
		cfg.Concurrency = fileCfg.Concurrency
	}
	if !cfg.DisableBuiltin {
		cfg.DisableBuiltin = fileCfg.DisableBuiltin
	}
	if cfg.MaxSize == 10*1024*1024 && fileCfg.MaxSize != "" {
		if v := parseMaxSize(fileCfg.MaxSize); v > 0 {
			cfg.MaxSize = v
		}
	}
	if !cfg.NoColor {
		cfg.NoColor = fileCfg.NoColor
	}
	if cfg.OutputFile == "" {
		cfg.OutputFile = fileCfg.OutputFile
	}
	if !cfg.Summary {
		cfg.Summary = fileCfg.Summary
	}
	if !cfg.AIEnabled {
		cfg.AIEnabled = fileCfg.AIEnabled
	}
	if cfg.AIModel == "" {
		cfg.AIModel = fileCfg.AIModel
	}
	if cfg.AIEndpoint == "" {
		cfg.AIEndpoint = fileCfg.AIEndpoint
	}
	return cfg
}

func parseMaxSize(s string) int64 {
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
		return 0
	}
	return v * multiplier
}
