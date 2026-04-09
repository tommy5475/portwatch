package daemon

import (
	"sync"
	"time"
)

// rateLimiter throttles alert delivery so that a burst of port changes
// does not flood downstream notification channels.
type rateLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	max      int
	tokens   int
	last     time.Time
}

// newRateLimiter returns a rateLimiter that allows at most max events per
// interval. max must be >= 1 and interval must be > 0; otherwise defaults
// of 10 events / 1 minute are used.
func newRateLimiter(max int, interval time.Duration) *rateLimiter {
	if max < 1 {
		max = 10
	}
	if interval <= 0 {
		interval = time.Minute
	}
	return &rateLimiter{
		interval: interval,
		max:      max,
		tokens:   max,
		last:     time.Now(),
	}
}

// Allow returns true if the event is permitted under the current rate limit.
// Tokens are replenished proportionally to the time elapsed since the last
// call, up to the configured maximum.
func (r *rateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.last)
	if elapsed >= r.interval {
		r.tokens = r.max
		r.last = now
	} else {
		// Partial replenishment.
		added := int(float64(r.max) * elapsed.Seconds() / r.interval.Seconds())
		if added > 0 {
			r.tokens += added
			if r.tokens > r.max {
				r.tokens = r.max
			}
			r.last = now
		}
	}

	if r.tokens <= 0 {
		return false
	}
	r.tokens--
	return true
}

// Remaining returns the number of tokens currently available.
func (r *rateLimiter) Remaining() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tokens
}
