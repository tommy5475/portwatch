package daemon

import (
	"sync"
	"sync/atomic"
	"time"
)

// fence is a one-shot write barrier that allows concurrent readers to wait
// until a single write-side event has been committed. Once crossed, the fence
// is permanently open and all subsequent Wait calls return immediately.
type fence struct {
	mu      sync.Mutex
	ready   chan struct{}
	crossed atomic.Bool
	crossAt time.Time
}

func newFence() *fence {
	return &fence{
		ready: make(chan struct{}),
	}
}

// Cross signals all waiting goroutines that the barrier has been passed.
// Subsequent calls to Cross are no-ops.
func (f *fence) Cross() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.crossed.Load() {
		return
	}
	f.crossAt = time.Now()
	f.crossed.Store(true)
	close(f.ready)
}

// Wait blocks until Cross has been called or ctx is done.
// Returns true if the fence was crossed, false if the context expired.
func (f *fence) Wait(ctx interface{ Done() <-chan struct{} }) bool {
	if f.crossed.Load() {
		return true
	}
	select {
	case <-f.ready:
		return true
	case <-ctx.Done():
		return false
	}
}

// Crossed reports whether Cross has been called.
func (f *fence) Crossed() bool {
	return f.crossed.Load()
}

// CrossedAt returns the time at which Cross was called, or the zero time
// if the fence has not yet been crossed.
func (f *fence) CrossedAt() time.Time {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.crossAt
}
