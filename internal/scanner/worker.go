package scanner

import (
	"sync"

	"github.com/wjames2000/skill-guard/internal/engine"
	pkgtypes "github.com/wjames2000/skill-guard/pkg/types"
)

func worker(
	jobs <-chan *pkgtypes.FileTarget,
	results chan<- []*pkgtypes.MatchResult,
	eng *engine.Engine,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	for target := range jobs {
		matches := eng.Match(target)
		if len(matches) > 0 {
			results <- matches
		}
	}
}
