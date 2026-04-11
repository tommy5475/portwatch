package daemon

import (
	"sync"
	"time"
)

// dedup suppresses duplicate events within a sliding time window.
// If the same key is seen more than once within the window duration,
// subsequent occurrences are dropped until the window expires.
type dedup struct {
	mu       sync.Mutex
	seen     map[string]time.Time
	window   time.Duration
	nowFn    func() time.Time
}

func newDedup(window time.Duration) *dedup {
	if window <= 0 {
		window = 5 * time.Second
	}
	return &dedup{
		seen:   make(map[string]time.Time),
		window: window,
		nowFn:  time.Now,
	}
}

// IsDuplicate reports whether key has been seen within the dedup window.
// If it has not been seen (or the window has expired), it records the key
// and returns false. Otherwise it returns true.
func (d *dedup) IsDuplicate(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFn()
	if t, ok := d.seen[key]; ok && now.Sub(t) < d.window {
		return true
	}
	d.seen[key] = now
	return false
}

// Evict removes all keys whose window has expired.
func (d *dedup) Evict() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFn()
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}

// Len returns the number of tracked keys.
func (d *dedup) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

// Reset clears all tracked keys.
func (d *dedup) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}
