package scanner

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wjames2000/skill-guard/internal/ai"
	"github.com/wjames2000/skill-guard/internal/engine"
	"github.com/wjames2000/skill-guard/internal/file"
	"github.com/wjames2000/skill-guard/internal/report"
	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func Scan(cfg *pkgtypes.Config) (*pkgtypes.ScanReport, error) {
	start := time.Now()

	files, err := file.Discover(cfg.Paths, &file.DiscoverOpts{
		Ignore:            cfg.Ignore,
		ExtInclude:        cfg.ExtInclude,
		ExtExclude:        cfg.ExtExclude,
		MaxSize:           cfg.MaxSize,
		Verbose:           cfg.Verbose,
		DiscoverGitIgnore: cfg.DiscoverGitIgnore,
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return emptyReport(start), nil
	}

	eng, err := engine.New(cfg.RulesFile, cfg.DisableBuiltin)
	if err != nil {
		return nil, err
	}

	numWorkers := cfg.Concurrency
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	jobs := make(chan *pkgtypes.FileTarget, numWorkers)
	results := make(chan []*pkgtypes.MatchResult, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, eng, &wg)
	}

	go func() {
		for _, f := range files {
			jobs <- f
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var scanned atomic.Int64
	var allResults []*pkgtypes.MatchResult
	for res := range results {
		allResults = append(allResults, res...)
		if cfg.Verbose {
			n := scanned.Add(1)
			fmt.Fprintf(os.Stderr, "\r进度: %d / %d 文件", n, len(files))
		}
	}
	if cfg.Verbose {
		fmt.Fprintln(os.Stderr)
	}

	r := report.Build(allResults, len(files), start)
	return r, nil
}

func AIVerifyResults(results []*pkgtypes.MatchResult, cfg *pkgtypes.Config) []*pkgtypes.MatchResult {
	if !cfg.AIEnabled {
		return results
	}

	aiCfg := &ai.AIConfig{
		Enabled: cfg.AIEnabled,
		Model:   cfg.AIModel,
	}

	if !ai.IsAvailable("") {
		fmt.Fprintln(os.Stderr, "⚠️ AI 检测不可用: Ollama 未运行，跳过 AI 验证")
		return results
	}

	var verified []*pkgtypes.MatchResult
	for _, r := range results {
		if r.LineContent == "" {
			verified = append(verified, r)
			continue
		}
		malicious, reason, err := ai.Analyze(r.LineContent, aiCfg)
		if err != nil {
			verified = append(verified, r)
			continue
		}
		if malicious {
			verified = append(verified, r)
		} else if cfg.Verbose {
			fmt.Fprintf(os.Stderr, "AI 判定非恶意: [%s] %s (%s)\n", r.RuleID, r.FilePath, reason)
		}
	}
	return verified
}

func emptyReport(start time.Time) *pkgtypes.ScanReport {
	return &pkgtypes.ScanReport{
		ScanTime:    start.Format("2006-01-02T15:04:05-07:00"),
		Duration:    time.Since(start).String(),
		TotalFiles:  0,
		TotalIssues: 0,
		Results:     []*pkgtypes.MatchResult{},
		Summary:     &pkgtypes.Summary{},
	}
}
