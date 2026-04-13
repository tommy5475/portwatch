package daemon

import (
	"sync"
	"time"
)

// epoch tracks a monotonically increasing generation counter with timestamps.
// Each time advance is called the counter increments and the timestamp is
// recorded. Callers can use the epoch to detect whether state has changed
// between two observations.
type epoch struct {
	mu      sync.RWMutex
	gen     uint64
	started time.Time
	last    time.Time
}

func newEpoch() *epoch {
	now := time.Now()
	return &epoch{
		gen:     0,
		started: now,
		last:    now,
	}
}

// Advance increments the generation counter and records the current time.
func (e *epoch) advance() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.gen++
	e.last = time.Now()
}

// Generation returns the current generation counter.
func (e *epoch) generation() uint64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.gen
}

// Last returns the time of the most recent advance, or the creation time if
// advance has never been called.
func (e *epoch) lastAdvance() time.Time {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.last
}

// Age returns the duration since the last advance.
func (e *epoch) age() time.Duration {
	return time.Since(e.lastAdvance())
}
