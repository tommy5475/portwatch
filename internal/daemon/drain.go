package daemon

import (
	"sync"
	"sync/atomic"
	"time"
)

// drain coordinates a graceful-shutdown sequence: callers acquire a slot
// before doing work and release it when done; the owner calls Wait to block
// until all in-flight slots are released or the deadline expires.
type drain struct {
	mu       sync.Mutex
	inflight int64
	done     chan struct{}
	closed   atomic.Bool
}

func newDrain() *drain {
	return &drain{done: make(chan struct{})}
}

// Acquire registers one unit of in-flight work.
// Returns false if the drain has already been closed.
func (d *drain) Acquire() bool {
	if d.closed.Load() {
		return false
	}
	atomic.AddInt64(&d.inflight, 1)
	// Re-check after increment to handle a race with Close.
	if d.closed.Load() {
		atomic.AddInt64(&d.inflight, -1)
		return false
	}
	return true
}

// Release signals that one unit of in-flight work has completed.
func (d *drain) Release() {
	if atomic.AddInt64(&d.inflight, -1) == 0 && d.closed.Load() {
		d.mu.Lock()
		select {
		case <-d.done:
		default:
			close(d.done)
		}
		d.mu.Unlock()
	}
}

// Close marks the drain as closed so no new work may be acquired,
// then waits up to timeout for all in-flight work to finish.
// Returns true if all work completed before the deadline.
func (d *drain) Close(timeout time.Duration) bool {
	d.closed.Store(true)
	// If nothing is in-flight, close the done channel immediately.
	if atomic.LoadInt64(&d.inflight) == 0 {
		d.mu.Lock()
		select {
		case <-d.done:
		default:
			close(d.done)
		}
		d.mu.Unlock()
	}
	select {
	case <-d.done:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Inflight returns the current number of acquired-but-not-released slots.
func (d *drain) Inflight() int64 {
	return atomic.LoadInt64(&d.inflight)
}
