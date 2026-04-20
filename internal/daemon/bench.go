package daemon

import (
	"sync"
	"time"
)

// bench measures the duration of repeated operations and exposes
// aggregate statistics: count, total elapsed, min, max, and mean.
type bench struct {
	mu    sync.Mutex
	count int64
	total time.Duration
	min   time.Duration
	max   time.Duration
}

func newBench() *bench {
	return &bench{}
}

// record registers a single observed duration.
func (b *bench) record(d time.Duration) {
	if d < 0 {
		d = 0
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.count++
	b.total += d
	if b.count == 1 || d < b.min {
		b.min = d
	}
	if d > b.max {
		b.max = d
	}
}

// count returns the number of recorded observations.
func (b *bench) len() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.count
}

// mean returns the arithmetic mean duration, or zero if no observations.
func (b *bench) mean() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.count == 0 {
		return 0
	}
	return b.total / time.Duration(b.count)
}

// min returns the shortest recorded duration.
func (b *bench) minDuration() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.min
}

// max returns the longest recorded duration.
func (b *bench) maxDuration() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.max
}

// reset clears all recorded observations.
func (b *bench) reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.count = 0
	b.total = 0
	b.min = 0
	b.max = 0
}
