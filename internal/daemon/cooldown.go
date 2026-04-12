package daemon

import (
	"sync"
	"time"
)

// cooldown enforces a minimum quiet period between successive activations
// of a named action. Unlike throttle, cooldown tracks the last *completion*
// time rather than the last *start* time, making it suitable for actions
// that have variable duration.
type cooldown struct {
	mu       sync.Mutex
	entries  map[string]time.Time
	period   time.Duration
	nowFn    func() time.Time
}

func newCooldown(period time.Duration) *cooldown {
	if period <= 0 {
		period = 5 * time.Second
	}
	return &cooldown{
		entries: make(map[string]time.Time),
		period:  period,
		nowFn:   time.Now,
	}
}

// Ready reports whether the cooldown period has elapsed for the given key.
func (c *cooldown) Ready(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	last, ok := c.entries[key]
	if !ok {
		return true
	}
	return c.nowFn().Sub(last) >= c.period
}

// Mark records the current time as the completion time for key.
func (c *cooldown) Mark(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = c.nowFn()
}

// Reset removes the cooldown record for key, allowing it to fire immediately.
func (c *cooldown) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Remaining returns how long until the cooldown expires for key.
// Returns zero if the key is already ready.
func (c *cooldown) Remaining(key string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	last, ok := c.entries[key]
	if !ok {
		return 0
	}
	elapsed := c.nowFn().Sub(last)
	if elapsed >= c.period {
		return 0
	}
	return c.period - elapsed
}
