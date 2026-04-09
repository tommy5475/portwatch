package daemon

import (
	"sync"
	"time"
)

// cache is a simple TTL-based in-memory key/value store used to avoid
// redundant work between scan cycles (e.g. deduplicating alerts for
// ports that have already been reported within the current window).
type cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
	ttl     time.Duration
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

// newCache returns a cache whose entries expire after ttl.
// If ttl is <= 0 it defaults to 60 seconds.
func newCache(ttl time.Duration) *cache {
	if ttl <= 0 {
		ttl = 60 * time.Second
	}
	return &cache{
		entries: make(map[string]cacheEntry),
		ttl:     ttl,
	}
}

// Set stores value under key, overwriting any existing entry.
func (c *cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Get returns the value stored under key and whether it was found and
// has not yet expired. Expired entries are deleted on access.
func (c *cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(e.expiresAt) {
		delete(c.entries, key)
		return nil, false
	}
	return e.value, true
}

// Delete removes the entry for key if present.
func (c *cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Len returns the number of non-expired entries currently held.
func (c *cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	count := 0
	for k, e := range c.entries {
		if now.After(e.expiresAt) {
			delete(c.entries, k)
		} else {
			count++
		}
	}
	return count
}
