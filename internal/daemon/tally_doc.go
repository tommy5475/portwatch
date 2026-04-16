// Package daemon provides internal daemon primitives for portwatch.
//
// # Tally
//
// tally is a thread-safe named counter map with an optional ceiling value.
// It is useful for tracking per-key event counts (e.g., per-port scan hits,
// alert counts, or error frequencies) without unbounded growth.
//
// Usage:
//
//	t := newTally(100)   // ceiling of 100
//	t.Inc("tcp:8080")
//	t.Add("tcp:8080", 5)
//	v := t.Get("tcp:8080")
//	snap := t.Snapshot()
//	t.Reset("tcp:8080")
//
// All methods are safe for concurrent use.
package daemon
