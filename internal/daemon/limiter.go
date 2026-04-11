package daemon

import (
	"sync"
	"time"
)

// limiter enforces a maximum number of alert dispatches per time window.
// It is safe for concurrent use.
type limiter struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	buckets  map[string][]time.Time
	nowFn    func() time.Time
}

func newLimiter(max int, window time.Duration) *limiter {
	if max <= 0 {
		max = 10
	}
	if window <= 0 {
		window = time.Minute
	}
	return &limiter{
		max:     max,
		window:  window,
		buckets: make(map[string][]time.Time),
		nowFn:   time.Now,
	}
}

// Allow reports whether the caller identified by key is permitted to proceed.
// It evicts timestamps outside the current window before checking the count.
func (l *limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.nowFn()
	cutoff := now.Add(-l.window)

	times := l.buckets[key]
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.max {
		l.buckets[key] = valid
		return false
	}

	l.buckets[key] = append(valid, now)
	return true
}

// Reset clears the event history for the given key.
func (l *limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Remaining returns how many more events key may produce within the current window.
func (l *limiter) Remaining(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.nowFn()
	cutoff := now.Add(-l.window)

	count := 0
	for _, t := range l.buckets[key] {
		if t.After(cutoff) {
			count++
		}
	}
	return l.max - count
}
