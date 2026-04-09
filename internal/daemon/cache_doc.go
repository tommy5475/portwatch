// Package daemon contains the runtime components that drive the portwatch
// daemon loop.
//
// # Cache
//
// cache is a lightweight, goroutine-safe TTL store used to deduplicate
// work across consecutive scan cycles.
//
// Typical usage:
//
//	c := newCache(30 * time.Second)
//	c.Set("tcp:8080", true)
//	if _, ok := c.Get("tcp:8080"); ok {
//		// already seen within TTL window — skip
//	}
//
// Entries are lazily evicted: expired items are removed the first time
// they are accessed via Get or counted via Len. There is no background
// goroutine, keeping the implementation allocation-free between accesses.
package daemon
