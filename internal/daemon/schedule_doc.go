// Package daemon provides the core runtime components of portwatch.
//
// # Schedule
//
// schedule implements a variable-interval tick strategy used by the
// daemon's run-loop to adapt polling frequency based on scan health.
//
// On each tick the caller reports whether the previous scan succeeded
// or failed. A failure causes the interval to grow by a configurable
// multiplicative factor (default 1.5×) up to a configured ceiling.
// A success immediately restores the nominal interval.
//
// This prevents thundering-herd behaviour when the host is under load
// while still recovering quickly once conditions improve.
//
// Usage:
//
//	s, err := newSchedule(30*time.Second, 5*time.Minute, 2.0)
//	if err != nil { ... }
//
//	for {
//		interval := s.next(lastScanFailed)
//		time.Sleep(interval)
//		lastScanFailed = runScan()
//	}
package daemon
