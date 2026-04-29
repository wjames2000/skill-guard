package types

type MatchResult struct {
	RuleID      string `json:"rule_id"`
	RuleName    string `json:"rule_name"`
	Severity    string `json:"severity"`
	FilePath    string `json:"file_path"`
	LineNumber  int    `json:"line_number"`
	LineContent string `json:"line_content"`
	MatchType   string `json:"match_type"`
}
