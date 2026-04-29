package engine

import (
	"testing"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func TestNewRule_ValidPattern(t *testing.T) {
	r, err := NewRule(&pkgtypes.Rule{
		ID: "TEST-001", Pattern: `AKIA[A-Z0-9]{16}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("Rule 不应为 nil")
	}
}

func TestNewRule_InvalidPattern(t *testing.T) {
	_, err := NewRule(&pkgtypes.Rule{
		ID: "TEST-002", Pattern: `[invalid`,
	})
	if err == nil {
		t.Error("无效正则应返回错误")
	}
}

func TestMatchesFileType(t *testing.T) {
	r, _ := NewRule(&pkgtypes.Rule{
		ID: "TEST-003", Pattern: `test`, FileTypes: []string{".py", ".sh"},
	})
	tests := []struct {
		ext  string
		want bool
	}{
		{".py", true}, {".sh", true}, {".md", false}, {".PY", true},
	}
	for _, tt := range tests {
		if got := r.MatchesFileType(tt.ext); got != tt.want {
			t.Errorf("MatchesFileType(%q) = %v, want %v", tt.ext, got, tt.want)
		}
	}
	r2, _ := NewRule(&pkgtypes.Rule{ID: "TEST-004", Pattern: `test`})
	if !r2.MatchesFileType(".anything") {
		t.Error("未限定 FileTypes 应匹配所有")
	}
}
