// Package daemon provides internal runtime primitives for the portwatch daemon.
//
// # Roster
//
// A roster is a thread-safe registry of named workers. Each worker can:
//
//   - join  – register itself with a timestamp
//   - checkin – update its last-seen time to signal liveness
//   - leave – deregister itself when it stops
//
// The roster also supports a markStale sweep that transitions workers to
// an "not alive" state when they have not checked in within a given TTL.
// This is useful for detecting hung or crashed goroutines without requiring
// an explicit leave call.
//
// Usage:
//
//	r := newRoster()
//	r.join("scanner-1")
//	r.checkin("scanner-1")
//	stale := r.markStale(30 * time.Second)
package daemon
