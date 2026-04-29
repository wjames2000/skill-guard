package file

import (
	"os"
	"path/filepath"
	"strings"
)

var DefaultIgnorePatterns = []string{
	".git", "node_modules", "vendor",
	"__pycache__", ".svn", ".hg", ".idea", ".vscode",
}

var DefaultExtInclude = []string{
	".md", ".json", ".yaml", ".yml",
	".py", ".sh", ".js", ".ts",
	".toml", ".xml", ".txt",
	".cfg", ".conf", ".ini", ".env",
	".bat", ".ps1", ".rb", ".php", ".lua",
}

func matchIgnore(path string, patterns []string) bool {
	base := filepath.Base(path)
	for _, p := range patterns {
		if matched, _ := filepath.Match(p, base); matched {
			return true
		}
		if strings.Contains(path, "/"+p+"/") || strings.HasSuffix(path, "/"+p) {
			return true
		}
	}
	return false
}

func checkExtension(ext string, include, exclude []string) bool {
	if ext == "" {
		return false
	}
	for _, e := range exclude {
		if strings.EqualFold(ext, e) {
			return false
		}
	}
	if len(include) > 0 {
		for _, e := range include {
			if strings.EqualFold(ext, e) {
				return true
			}
		}
		return false
	}
	return true
}

func LoadGitIgnore(dir string) []string {
	path := filepath.Join(dir, ".gitignore")
	data, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return nil
	}
	var patterns []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "!") {
			continue
		}
		patterns = append(patterns, strings.TrimSuffix(line, "/"))
	}
	return patterns
}

func isDefaultExt(ext string) bool {
	for _, e := range DefaultExtInclude {
		if ext == e {
			return true
		}
	}
	return false
}
