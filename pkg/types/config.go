package types

import (
	"errors"
	"strings"
)

var validSeverities = map[string]bool{
	"critical": true,
	"high":     true,
	"medium":   true,
	"low":      true,
}

type Config struct {
	Paths          []string
	ConfigFile     string
	RulesFile      string
	Severity       string
	JSONOutput     bool
	Quiet          bool
	Verbose        bool
	Ignore         []string
	ExtInclude     []string
	ExtExclude     []string
	MaxSize        int64
	Concurrency    int
	DisableBuiltin bool
}

func DefaultConfig() *Config {
	return &Config{
		Paths:       []string{"."},
		MaxSize:     10 * 1024 * 1024,
		Concurrency: 0,
	}
}

func (c *Config) Validate() error {
	if len(c.Paths) == 0 {
		return errors.New("至少需要一个扫描路径")
	}
	if c.Severity != "" && !validSeverities[strings.ToLower(c.Severity)] {
		return errors.New("无效的严重级别，可选: critical/high/medium/low")
	}
	if c.MaxSize <= 0 {
		return errors.New("文件大小上限必须大于 0")
	}
	return nil
}
