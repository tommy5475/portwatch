package daemon

import (
	"math"
	"time"
)

// backoff implements an exponential back-off strategy used when consecutive
// scan errors occur. The delay doubles on each failure up to MaxDelay, then
// resets to BaseDelay once a successful scan is recorded.
type backoff struct {
	BaseDelay time.Duration
	MaxDelay  time.Duration
	multiplier float64
	current   time.Duration
	failures  int
}

// newBackoff returns a backoff with sensible defaults.
// baseDelay is the initial wait; maxDelay caps the growth.
func newBackoff(baseDelay, maxDelay time.Duration) *backoff {
	if baseDelay <= 0 {
		baseDelay = 2 * time.Second
	}
	if maxDelay <= 0 || maxDelay < baseDelay {
		maxDelay = 60 * time.Second
	}
	return &backoff{
		BaseDelay:  baseDelay,
		MaxDelay:   maxDelay,
		multiplier: 2.0,
		current:    baseDelay,
	}
}

// Failure records a failed attempt and returns the delay that should be
// observed before the next retry.
func (b *backoff) Failure() time.Duration {
	b.failures++
	delay := time.Duration(float64(b.BaseDelay) * math.Pow(b.multiplier, float64(b.failures-1)))
	if delay > b.MaxDelay {
		delay = b.MaxDelay
	}
	b.current = delay
	return delay
}

// Success resets the back-off state after a successful attempt.
func (b *backoff) Success() {
	b.failures = 0
	b.current = b.BaseDelay
}

// Current returns the most recently computed delay without advancing state.
func (b *backoff) Current() time.Duration {
	return b.current
}

// Failures returns the number of consecutive failures recorded.
func (b *backoff) Failures() int {
	return b.failures
}
