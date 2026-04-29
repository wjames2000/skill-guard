package scanner

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wjames2000/skill-guard/internal/engine"
	"github.com/wjames2000/skill-guard/internal/file"
	"github.com/wjames2000/skill-guard/internal/report"
	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func Scan(cfg *pkgtypes.Config) (*pkgtypes.ScanReport, error) {
	start := time.Now()

	files, err := file.Discover(cfg.Paths, &file.DiscoverOpts{
		Ignore:     cfg.Ignore,
		ExtInclude: cfg.ExtInclude,
		ExtExclude: cfg.ExtExclude,
		MaxSize:    cfg.MaxSize,
		Verbose:    cfg.Verbose,
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
			println("进度:", n, "/", len(files))
		}
	}
	if cfg.Verbose {
		println("扫描完成")
	}

	r := report.Build(allResults, len(files), start)
	return r, nil
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
