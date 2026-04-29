package engine

import pkgtypes "github.com/hpds.cc/skill-guard/pkg/types"

func BuiltinRules() []*pkgtypes.Rule {
	return []*pkgtypes.Rule{
		{
			ID: "SKL-001", Name: "硬编码 AWS Access Key",
			Severity: "Critical", Pattern: `(?i)AKIA[A-Z0-9]{16}`,
			Keywords: []string{"AKIA"},
		},
	}
}
