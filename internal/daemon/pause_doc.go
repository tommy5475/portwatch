// Package daemon provides the core runtime primitives for the portwatch
// daemon process.
//
// # Pause
//
// pause is a lightweight pausable-timer helper that lets the scan loop be
// temporarily suspended (e.g. during a maintenance window or after the
// circuit breaker opens) and later resumed without drifting the nominal
// scan interval.
//
// Typical usage:
//
//	p := newPause()
//
//	// Suspend scanning.
//	p.Pause()
//
//	// ... wait for condition ...
//
//	// Resume scanning.
//	p.Resume()
//
//	// Inspect how long we were idle.
//	fmt.Println(p.TotalPaused())
//
// All methods are safe for concurrent use.
package daemon
