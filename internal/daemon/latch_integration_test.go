package daemon

import (
	"sync"
	"sync/atomic"
	"testing"
)

// TestLatchGatesWorkUntilReady simulates a readiness gate pattern:
// workers wait until the latch is set before counting their work.
func TestLatchGatesWorkUntilReady(t *testing.T) {
	l := newLatch()
	var (
		wg      sync.WaitGroup
		counted int64
	)

	const workers = 20
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.IfSet(func() {
				atomic.AddInt64(&counted, 1)
			})
		}()
	}

	// None should have run yet.
	wg.Wait()
	if counted != 0 {
		t.Fatalf("expected 0 counted before Set, got %d", counted)
	}

	// Now set and re-run.
	l.Set()
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.IfSet(func() {
				atomic.AddInt64(&counted, 1)
			})
		}()
	}
	wg.Wait()
	if counted != workers {
		t.Fatalf("expected %d counted after Set, got %d", workers, counted)
	}
}
