package output

import (
	"encoding/json"
	"fmt"
	"io"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

type sarifLog struct {
	Version string    `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string       `json:"name"`
	Version        string       `json:"version"`
	InformationURI string       `json:"informationUri"`
	Rules          []sarifRule  `json:"rules"`
}

type sarifRule struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	ShortDesc   sarifMessage `json:"shortDescription"`
	FullDesc    sarifMessage `json:"fullDescription"`
	DefaultConf sarifConfig  `json:"defaultConfiguration"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifConfig struct {
	Level string `json:"level"`
}

type sarifResult struct {
	RuleID    string         `json:"ruleId"`
	Message   sarifMessage   `json:"message"`
	Level     string         `json:"level"`
	Locations []sarifLocation `json:"locations"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine int `json:"startLine"`
}

type SARIFRenderer struct{}

func (s *SARIFRenderer) Render(w io.Writer, report *pkgtypes.ScanReport) error {
	var rules []sarifRule
	for _, r := range report.Results {
		rules = append(rules, sarifRule{
			ID:   r.RuleID,
			Name: r.RuleName,
			ShortDesc: sarifMessage{Text: fmt.Sprintf("%s: %s", r.Severity, r.RuleName)},
			FullDesc:  sarifMessage{Text: r.RuleName},
			DefaultConf: sarifConfig{Level: toSARIFLevel(r.Severity)},
		})
	}
	rules = dedupSARIFRules(rules)

	var results []sarifResult
	for _, r := range report.Results {
		results = append(results, sarifResult{
			RuleID:  r.RuleID,
			Message: sarifMessage{Text: r.LineContent},
			Level:   toSARIFLevel(r.Severity),
			Locations: []sarifLocation{{
				PhysicalLocation: sarifPhysicalLocation{
					ArtifactLocation: sarifArtifactLocation{URI: r.FilePath},
					Region:           sarifRegion{StartLine: r.LineNumber},
				},
			}},
		})
	}

	log := sarifLog{
		Version: "2.1.0",
		Runs: []sarifRun{{
			Tool: sarifTool{
				Driver: sarifDriver{
					Name:           "skill-guard",
					Version:        "0.1.0",
					InformationURI: "https://github.com/wjames2000/skill-guard",
					Rules:          rules,
				},
			},
			Results: results,
		}},
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(log)
}

func toSARIFLevel(severity string) string {
	switch severity {
	case "Critical", "High":
		return "error"
	case "Medium":
		return "warning"
	case "Low":
		return "note"
	default:
		return "none"
	}
}

func dedupSARIFRules(rules []sarifRule) []sarifRule {
	seen := make(map[string]bool)
	var unique []sarifRule
	for _, r := range rules {
		if seen[r.ID] {
			continue
		}
		seen[r.ID] = true
		unique = append(unique, r)
	}
	return unique
}
