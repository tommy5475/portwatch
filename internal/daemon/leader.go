package daemon

import (
	"sync"
	"time"
)

// leader implements a simple single-process leadership latch that tracks
// whether this daemon instance considers itself the active leader. It is
// useful when multiple portwatch processes share a state file and only one
// should emit alerts at a time.
type leader struct {
	mu        sync.RWMutex
	active    bool
	acquiredAt time.Time
	renewedAt  time.Time
	term      uint64
}

func newLeader() *leader {
	return &leader{}
}

// Acquire marks this instance as the active leader. It is idempotent; calling
// Acquire when already leading increments the term and updates renewedAt.
func (l *leader) Acquire() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if !l.active {
		l.active = true
		l.acquiredAt = now
	}
	l.renewedAt = now
	l.term++
}

// Release relinquishes leadership. Safe to call when not leading.
func (l *leader) Release() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.active = false
}

// IsLeader reports whether this instance currently holds leadership.
func (l *leader) IsLeader() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.active
}

// Term returns the monotonically increasing leadership term counter.
func (l *leader) Term() uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.term
}

// Age returns how long this instance has held leadership continuously.
// Returns zero if not currently the leader.
func (l *leader) Age() time.Duration {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if !l.active {
		return 0
	}
	return time.Since(l.acquiredAt)
}

// LastRenewed returns the time of the most recent Acquire call, or the zero
// time if leadership has never been acquired.
func (l *leader) LastRenewed() time.Time {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.renewedAt
}
