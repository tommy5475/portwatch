package daemon

import (
	"sync"
	"time"
)

// token is a one-shot cancellable ticket that expires after a fixed TTL.
// It is useful for tracking in-flight work items that must be acknowledged
// within a deadline; if the token is not claimed before it expires it is
// automatically marked as timed-out.
type token struct {
	mu        sync.Mutex
	id        string
	issuedAt  time.Time
	ttl       time.Duration
	claimed   bool
	timedOut  bool
	stopTimer func()
}

func newToken(id string, ttl time.Duration, onExpire func(id string)) *token {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	t := &token{
		id:       id,
		issuedAt: time.Now(),
		ttl:      ttl,
	}
	timer := time.AfterFunc(ttl, func() {
		t.mu.Lock()
		if !t.claimed {
			t.timedOut = true
		}
		t.mu.Unlock()
		if onExpire != nil {
			onExpire(id)
		}
	})
	t.stopTimer = func() { timer.Stop() }
	return t
}

// Claim marks the token as successfully consumed. Returns false if the token
// was already claimed or has already timed out.
func (t *token) Claim() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.claimed || t.timedOut {
		return false
	}
	t.claimed = true
	t.stopTimer()
	return true
}

// ID returns the token identifier.
func (t *token) ID() string { return t.id }

// IsClaimed reports whether the token was successfully claimed.
func (t *token) IsClaimed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.claimed
}

// IsTimedOut reports whether the token expired before being claimed.
func (t *token) IsTimedOut() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.timedOut
}

// Age returns how long ago the token was issued.
func (t *token) Age() time.Duration {
	return time.Since(t.issuedAt)
}
