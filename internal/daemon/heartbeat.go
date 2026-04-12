package daemon

import (
	"sync"
	"time"
)

// heartbeat tracks periodic liveness pulses emitted by the scan loop.
// It exposes the last beat time and elapsed duration since the last beat,
// which the health endpoint and watchdog can query without blocking.
type heartbeat struct {
	mu       sync.RWMutex
	lastBeat time.Time
	total    int64
	started  time.Time
}

func newHeartbeat() *heartbeat {
	return &heartbeat{
		started: time.Now(),
	}
}

// pulse records a new heartbeat at the current wall-clock time.
func (h *heartbeat) pulse() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lastBeat = time.Now()
	h.total++
}

// last returns the time of the most recent pulse.
// A zero time indicates no pulse has been recorded yet.
func (h *heartbeat) last() time.Time {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.lastBeat
}

// age returns how long ago the last pulse occurred.
// Returns the full uptime if no pulse has been recorded.
func (h *heartbeat) age() time.Duration {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if h.lastBeat.IsZero() {
		return time.Since(h.started)
	}
	return time.Since(h.lastBeat)
}

// count returns the total number of pulses recorded since creation.
func (h *heartbeat) count() int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.total
}

// alive reports whether a pulse was recorded within the given threshold.
func (h *heartbeat) alive(threshold time.Duration) bool {
	return h.age() <= threshold
}
