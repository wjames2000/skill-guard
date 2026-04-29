package config

type FileConfig struct {
	Ignore         []string `yaml:"ignore,omitempty"`
	Rules          string   `yaml:"rules,omitempty"`
	Severity       string   `yaml:"severity,omitempty"`
	ExtInclude     []string `yaml:"ext_include,omitempty"`
	ExtExclude     []string `yaml:"ext_exclude,omitempty"`
	MaxSize        string   `yaml:"max_size,omitempty"`
	Concurrency    int      `yaml:"concurrency,omitempty"`
	DisableBuiltin bool     `yaml:"disable_builtin,omitempty"`
}
