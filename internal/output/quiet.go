package output

import (
	"fmt"
	"io"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

type QuietRenderer struct{}

func (q *QuietRenderer) Render(w io.Writer, report *pkgtypes.ScanReport) error {
	if len(report.Results) == 0 {
		return nil
	}
	seen := make(map[string]bool)
	for _, r := range report.Results {
		if seen[r.FilePath] {
			continue
		}
		seen[r.FilePath] = true
		fmt.Fprintln(w, r.FilePath)
	}
	return nil
}
