package daemon

import (
	"sync"
	"time"
)

// throttle suppresses repeated alert dispatches for the same port within a
// configurable cooldown window. This prevents alert storms when a port flaps.
type throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

func newThrottle(cooldown time.Duration) *throttle {
	if cooldown <= 0 {
		cooldown = 30 * time.Second
	}
	return &throttle{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true when the given key has not been seen within the cooldown
// window. Calling Allow records the current time for the key.
func (t *throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.last[key]; ok && now.Sub(last) < t.cooldown {
		return false
	}
	t.last[key] = now
	return true
}

// Reset clears the recorded timestamp for key, allowing the next call to
// Allow to succeed immediately regardless of the cooldown window.
func (t *throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Len returns the number of tracked keys.
func (t *throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.last)
}
