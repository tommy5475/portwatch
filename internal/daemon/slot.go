package daemon

import (
	"sync"
	"time"
)

// slot is a named reservation that can be acquired for a fixed duration.
// It is useful for time-boxing exclusive access to a shared resource.
type slot struct {
	mu       sync.Mutex
	name     string
	held     bool
	acquiredAt time.Time
	ttl      time.Duration
	count    int64
}

func newSlot(name string, ttl time.Duration) *slot {
	if ttl <= 0 {
		ttl = 5 * time.Second
	}
	if name == "" {
		name = "default"
	}
	return &slot{name: name, ttl: ttl}
}

// Acquire attempts to hold the slot. Returns true if successful.
// If the slot is already held and the TTL has not expired, it returns false.
func (s *slot) Acquire() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.held && time.Since(s.acquiredAt) < s.ttl {
		return false
	}
	s.held = true
	s.acquiredAt = time.Now()
	s.count++
	return true
}

// Release frees the slot so another caller may acquire it.
func (s *slot) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.held = false
}

// Held reports whether the slot is currently held and within its TTL.
func (s *slot) Held() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.held && time.Since(s.acquiredAt) < s.ttl
}

// Count returns the total number of times the slot has been acquired.
func (s *slot) Count() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.count
}

// Name returns the slot's name.
func (s *slot) Name() string {
	return s.name
}
