package daemon

import "sync"

// latch is a one-shot boolean flag that can be set once and read many times.
// Once set, it cannot be cleared. It is safe for concurrent use.
type latch struct {
	mu  sync.RWMutex
	set bool
}

// newLatch returns a new latch in the unset state.
func newLatch() *latch {
	return &latch{}
}

// Set marks the latch as triggered. Subsequent calls are no-ops.
func (l *latch) Set() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.set = true
}

// IsSet reports whether the latch has been set.
func (l *latch) IsSet() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.set
}

// IfSet calls fn exactly once if the latch is currently set.
// Returns true if fn was called.
func (l *latch) IfSet(fn func()) bool {
	if l.IsSet() {
		fn()
		return true
	}
	return false
}
