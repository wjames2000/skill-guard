package types

type Summary struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}

type ScanReport struct {
	ScanTime    string         `json:"scan_time"`
	Duration    string         `json:"duration"`
	TotalFiles  int            `json:"total_files"`
	TotalIssues int            `json:"total_issues"`
	Results     []*MatchResult `json:"results"`
	Summary     *Summary       `json:"summary"`
}
