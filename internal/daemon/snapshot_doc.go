// Package daemon provides the long-running daemon that ties together
// scanning, filtering, alerting, and reporting.
//
// # Snapshot
//
// The snapshot type is an internal, concurrency-safe cache of the most
// recently completed port scan result. It is updated by the scan loop
// after every successful scan and read by the health-check endpoint and
// metrics recorder without holding any long-lived locks.
//
// Usage:
//
//	snap := newSnapshot()
//
//	// writer side (scan loop)
//	snap.update(portMap)
//
//	// reader side (health endpoint, metrics)
//	ports, capturedAt := snap.get()
//	age := snap.age()
//	count := snap.count()
//
// All methods are safe for concurrent use.
package daemon
