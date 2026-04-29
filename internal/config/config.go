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
	NoColor        bool     `yaml:"no_color,omitempty"`
	OutputFile     string   `yaml:"output,omitempty"`
	Summary        bool     `yaml:"summary,omitempty"`
	AIEnabled      bool     `yaml:"ai_enabled,omitempty"`
	AIModel        string   `yaml:"ai_model,omitempty"`
	AIEndpoint     string   `yaml:"ai_endpoint,omitempty"`
}
