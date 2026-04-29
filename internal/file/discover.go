package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

type DiscoverOpts struct {
	Ignore            []string
	ExtInclude        []string
	ExtExclude        []string
	MaxSize           int64
	Verbose           bool
	DiscoverGitIgnore bool
}

func Discover(roots []string, opts *DiscoverOpts) ([]*pkgtypes.FileTarget, error) {
	if opts == nil {
		opts = &DiscoverOpts{}
	}

	ignorePatterns := append([]string{}, DefaultIgnorePatterns...)
	ignorePatterns = append(ignorePatterns, opts.Ignore...)

	if opts.DiscoverGitIgnore {
		for _, root := range roots {
			gitignorePatterns := LoadGitIgnore(root)
			ignorePatterns = append(ignorePatterns, gitignorePatterns...)
		}
	}

	var files []*pkgtypes.FileTarget
	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				if matchIgnore(path, ignorePatterns) {
					return filepath.SkipDir
				}
				return nil
			}
			if d.Type()&os.ModeSymlink != 0 {
				return nil
			}
			if matchIgnore(path, ignorePatterns) {
				return nil
			}

			info, err := d.Info()
			if err != nil {
				return nil
			}

			ext := strings.ToLower(filepath.Ext(path))
			if !checkExtension(ext, opts.ExtInclude, opts.ExtExclude) {
				return nil
			}
			if len(opts.ExtInclude) == 0 && len(opts.ExtExclude) == 0 {
				if !isDefaultExt(ext) {
					return nil
				}
			}
			if opts.MaxSize > 0 && info.Size() > opts.MaxSize {
				if opts.Verbose {
					fmt.Fprintf(os.Stderr, "跳过超限: %s (%d bytes)\n", path, info.Size())
				}
				return nil
			}

			relPath, _ := filepath.Rel(root, path)
			files = append(files, &pkgtypes.FileTarget{
				Path:    path,
				RelPath: relPath,
				Size:    info.Size(),
				Ext:     ext,
			})
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("遍历目录失败 %s: %w", root, err)
		}
	}
	return files, nil
}
