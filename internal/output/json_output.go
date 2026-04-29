package output

import (
	"encoding/json"
	"io"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

type JSONRenderer struct{}

func (j *JSONRenderer) Render(w io.Writer, report *pkgtypes.ScanReport) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
