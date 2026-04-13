package daemon

import (
	"sync"
	"sync/atomic"
	"time"
)

// epoch tracks a monotonically incrementing generation counter alongside the
// wall-clock time at which each generation began. It is safe for concurrent
// use.
type epoch struct {
	mu    sync.Mutex
	gen   atomic.Uint64
	start time.Time
}

func newEpoch() *epoch {
	e := &epoch{start: time.Now()}
	e.gen.Store(1)
	return e
}

// Advance increments the generation counter and records the current time as
// the start of the new epoch. It returns the new generation number.
func (e *epoch) Advance() uint64 {
	e.mu.Lock()
	e.start = time.Now()
	e.mu.Unlock()
	return e.gen.Add(1)
}

// Current returns the current generation number without advancing it.
func (e *epoch) Current() uint64 {
	return e.gen.Load()
}

// Age returns how long the current epoch has been active.
func (e *epoch) Age() time.Duration {
	e.mu.Lock()
	start := e.start
	e.mu.Unlock()
	return time.Since(start)
}

// Reset sets the generation back to 1 and records a fresh start time.
func (e *epoch) Reset() {
	e.mu.Lock()
	e.start = time.Now()
	e.mu.Unlock()
	e.gen.Store(1)
}
