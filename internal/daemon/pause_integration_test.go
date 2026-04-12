package daemon

import (
	"sync"
	"testing"
	"time"
)

// TestPauseConcurrentPauseResume verifies that rapid concurrent Pause/Resume
// calls do not race or panic.
func TestPauseConcurrentPauseResume(t *testing.T) {
	p := newPause()
	var wg sync.WaitGroup
	const workers = 20

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				p.Pause()
			} else {
				p.Resume()
			}
		}(i)
	}
	wg.Wait()
	// No assertion on final state — just must not race or panic.
}

// TestPauseTotalMonotonicallyIncreases verifies the cumulative total never
// decreases across multiple pause/resume cycles.
func TestPauseTotalMonotonicallyIncreases(t *testing.T) {
	p := newPause()
	var prev time.Duration

	for i := 0; i < 3; i++ {
		p.Pause()
		time.Sleep(15 * time.Millisecond)
		p.Resume()

		current := p.TotalPaused()
		if current < prev {
			t.Fatalf("cycle %d: TotalPaused decreased from %v to %v", i, prev, current)
		}
		prev = current
	}
}
