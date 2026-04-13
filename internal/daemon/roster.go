package daemon

import (
	"sync"
	"time"
)

// roster tracks a named set of active workers, recording when each
// registered, last checked in, and whether it is still alive.
type roster struct {
	mu      sync.RWMutex
	entries map[string]*rosterEntry
}

type rosterEntry struct {
	name      string
	joined    time.Time
	lastSeen  time.Time
	alive     bool
}

func newRoster() *roster {
	return &roster{entries: make(map[string]*rosterEntry)}
}

// join registers a worker by name. Calling join again resets the entry.
func (r *roster) join(name string) {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[name] = &rosterEntry{
		name:     name,
		joined:   now,
		lastSeen: now,
		alive:    true,
	}
}

// checkin updates the lastSeen timestamp for a worker.
// Returns false if the worker is not registered.
func (r *roster) checkin(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[name]
	if !ok {
		return false
	}
	e.lastSeen = time.Now()
	e.alive = true
	return true
}

// leave removes a worker from the roster.
func (r *roster) leave(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, name)
}

// alive reports whether the named worker is registered and marked alive.
func (r *roster) isAlive(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[name]
	return ok && e.alive
}

// len returns the number of currently registered workers.
func (r *roster) len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}

// snapshot returns a copy of all entries.
func (r *roster) snapshot() []rosterEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]rosterEntry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, *e)
	}
	return out
}

// markStale sets alive=false for any worker whose lastSeen is older than ttl.
func (r *roster) markStale(ttl time.Duration) int {
	cutoff := time.Now().Add(-ttl)
	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	for _, e := range r.entries {
		if e.alive && e.lastSeen.Before(cutoff) {
			e.alive = false
			count++
		}
	}
	return count
}
