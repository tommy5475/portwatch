package daemon

import (
	"sync"
	"time"
)

// pause is a pausable timer that allows the daemon scan loop to be
// temporarily suspended and resumed without losing its interval rhythm.
type pause struct {
	mu       sync.Mutex
	paused   bool
	pausedAt time.Time
	total    time.Duration
}

func newPause() *pause {
	return &pause{}
}

// Pause suspends the timer. Calling Pause when already paused is a no-op.
func (p *pause) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.paused {
		return
	}
	p.paused = true
	p.pausedAt = time.Now()
}

// Resume resumes the timer. Calling Resume when not paused is a no-op.
func (p *pause) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.paused {
		return
	}
	p.total += time.Since(p.pausedAt)
	p.paused = false
}

// IsPaused reports whether the timer is currently paused.
func (p *pause) IsPaused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.paused
}

// TotalPaused returns the cumulative duration spent paused.
// If currently paused, the ongoing pause is included.
func (p *pause) TotalPaused() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.paused {
		return p.total + time.Since(p.pausedAt)
	}
	return p.total
}

// Reset clears the cumulative pause duration and sets the state to unpaused.
func (p *pause) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = false
	p.total = 0
	p.pausedAt = time.Time{}
}
