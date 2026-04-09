// Package daemon contains the runtime engine for portwatch.
//
// # Throttle
//
// The throttle type prevents alert storms caused by rapidly flapping ports.
// When a port opens and closes repeatedly within a short window the alert
// subsystem would otherwise emit a notification on every scan cycle.
//
// Usage:
//
//	th := newThrottle(30 * time.Second)
//	if th.Allow("tcp:8080:opened") {
//		// dispatch alert
//	}
//
// Allow is safe for concurrent use. Reset can be used in tests or when an
// operator explicitly clears the suppression state for a specific key.
package daemon
