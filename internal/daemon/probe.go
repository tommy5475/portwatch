package daemon

import (
	"sync"
	"time"
)

// probe tracks the liveness and readiness state of a component.
// It is safe for concurrent use.
type probe struct {
	mu        sync.RWMutex
	live      bool
	ready     bool
	liveAt    time.Time
	readyAt   time.Time
	failedAt  time.Time
	failCount int
}

func newProbe() *probe {
	return &probe{live: true}
}

// MarkLive sets the component as live (healthy heartbeat).
func (p *probe) MarkLive() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.live = true
	p.liveAt = time.Now()
}

// MarkNotLive marks the component as not live and records the failure.
func (p *probe) MarkNotLive() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.live = false
	p.failCount++
	p.failedAt = time.Now()
}

// MarkReady signals the component is ready to serve traffic.
func (p *probe) MarkReady() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ready = true
	p.readyAt = time.Now()
}

// MarkNotReady signals the component is temporarily unavailable.
func (p *probe) MarkNotReady() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ready = false
}

// IsLive returns whether the component is currently live.
func (p *probe) IsLive() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.live
}

// IsReady returns whether the component is currently ready.
func (p *probe) IsReady() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.ready
}

// FailCount returns the total number of times MarkNotLive was called.
func (p *probe) FailCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.failCount
}

// FailedAt returns the time of the most recent liveness failure, or zero.
func (p *probe) FailedAt() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.failedAt
}

// Reset clears all state back to initial values.
func (p *probe) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.live = true
	p.ready = false
	p.liveAt = time.Time{}
	p.readyAt = time.Time{}
	p.failedAt = time.Time{}
	p.failCount = 0
}
