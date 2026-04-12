package daemon

import (
	"sync"
	"time"
)

// checkpoint records named milestones reached during daemon operation.
// Each checkpoint stores the time it was last reached and how many times
// it has been reached in total. Safe for concurrent use.
type checkpoint struct {
	mu      sync.RWMutex
	entries map[string]*checkpointEntry
}

type checkpointEntry struct {
	count   int64
	lastAt  time.Time
	firstAt time.Time
}

func newCheckpoint() *checkpoint {
	return &checkpoint{
		entries: make(map[string]*checkpointEntry),
	}
}

// Mark records that the named checkpoint has been reached.
func (c *checkpoint) Mark(name string) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[name]
	if !ok {
		c.entries[name] = &checkpointEntry{count: 1, firstAt: now, lastAt: now}
		return
	}
	e.count++
	e.lastAt = now
}

// Count returns how many times the named checkpoint has been reached.
func (c *checkpoint) Count(name string) int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if e, ok := c.entries[name]; ok {
		return e.count
	}
	return 0
}

// LastAt returns the time the named checkpoint was most recently reached.
// Returns zero time if the checkpoint has never been reached.
func (c *checkpoint) LastAt(name string) time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if e, ok := c.entries[name]; ok {
		return e.lastAt
	}
	return time.Time{}
}

// FirstAt returns the time the named checkpoint was first reached.
// Returns zero time if the checkpoint has never been reached.
func (c *checkpoint) FirstAt(name string) time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if e, ok := c.entries[name]; ok {
		return e.firstAt
	}
	return time.Time{}
}

// Reset clears all recorded checkpoints.
func (c *checkpoint) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*checkpointEntry)
}

// Names returns the names of all checkpoints that have been reached at least once.
func (c *checkpoint) Names() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	names := make([]string, 0, len(c.entries))
	for k := range c.entries {
		names = append(names, k)
	}
	return names
}
