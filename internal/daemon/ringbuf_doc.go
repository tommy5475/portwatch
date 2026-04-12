// Package daemon provides the core runtime components for the portwatch
// daemon, including scheduling, resilience primitives, and data structures.
//
// # Ring Buffer
//
// ringbuf[T] is a generic, fixed-capacity circular buffer suitable for
// maintaining a rolling history of events such as recent scan results,
// alert records, or error messages.
//
// When the buffer is full, new entries overwrite the oldest entry so that
// memory usage remains bounded regardless of how long the daemon runs.
//
// Usage:
//
//	buf := newRingbuf[string](32)
//	buf.Push("port 8080 opened")
//	entries := buf.Snapshot() // oldest-first slice copy
//
// All methods are safe for concurrent use.
package daemon
