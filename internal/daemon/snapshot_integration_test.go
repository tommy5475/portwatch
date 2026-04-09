package daemon

import (
	"sync"
	"testing"

	"portwatch/internal/state"
)

// TestSnapshotConcurrentReadWrite verifies that concurrent updates and
// reads do not race (run with -race to exercise the mutex paths).
func TestSnapshotConcurrentReadWrite(t *testing.T) {
	s := newSnapshot()

	const writers = 4
	const readers = 8
	const iterations = 200

	var wg sync.WaitGroup

	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				pm := state.PortMap{
					"tcp:80": {Port: 80, Protocol: "tcp", Open: true},
				}
				s.update(pm)
			}
		}(w)
	}

	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_, _ = s.get()
				_ = s.age()
				_ = s.count()
			}
		}()
	}

	wg.Wait()

	if s.count() != int64(writers*iterations) {
		t.Errorf("expected %d updates, got %d", writers*iterations, s.count())
	}
}
