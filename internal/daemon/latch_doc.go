// Package daemon provides internal runtime primitives for the portwatch daemon.
//
// # Latch
//
// A latch is a one-shot, write-once boolean flag. Once set it remains set
// for the lifetime of the value. It is designed for signalling irreversible
// state transitions such as "first scan completed" or "degraded mode entered".
//
// Example usage:
//
//	ready := newLatch()
//
//	// in scan goroutine:
//	ready.Set()
//
//	// in health-check handler:
//	if !ready.IsSet() {
//		http.Error(w, "not ready", http.StatusServiceUnavailable)
//		return
//	}
//
// Latch is safe for concurrent use by multiple goroutines.
package daemon
