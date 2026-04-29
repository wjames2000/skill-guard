package types

type Rule struct {
	ID          string   `yaml:"id" json:"id"`
	Name        string   `yaml:"name" json:"name"`
	Severity    string   `yaml:"severity" json:"severity"`
	Description string   `yaml:"description" json:"description"`
	Pattern     string   `yaml:"pattern" json:"pattern"`
	Keywords    []string `yaml:"keywords,omitempty" json:"keywords,omitempty"`
	FileTypes   []string `yaml:"file_types,omitempty" json:"file_types,omitempty"`
	Ref         string   `yaml:"ref,omitempty" json:"ref,omitempty"`
}

type RuleList struct {
	Rules []Rule `yaml:"rules" json:"rules"`
}
