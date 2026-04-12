package daemon

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestBarrierMultiRoundPipeline simulates a two-stage pipeline where workers
// must all complete stage-1 before any worker begins stage-2.
func TestBarrierMultiRoundPipeline(t *testing.T) {
	const workers = 8
	const rounds = 4

	b := newBarrier(workers)
	var stage1Done int64
	var stage2Started int64

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for r := 0; r < rounds; r++ {
				// Stage 1
				atomic.AddInt64(&stage1Done, 1)
				b.Wait()
			// After barrier: all stage-1 increments must be visible.
				if v := atomic.LoadInt64(&stage1Done); v < int64(workers*(r+1)) {
					t.Errorf("round %d: stage1Done=%d want >=%d", r, v, workers*(r+1))
				}
				// Stage 2
				atomic.AddInt64(&stage2Started, 1)
				b.Wait()
			}
		}()
	}
	wg.Wait()

	if want := int64(workers * rounds); atomic.LoadInt64(&stage2Started) != want {
		t.Fatalf("stage2Started=%d want %d", stage2Started, want)
	}
}
