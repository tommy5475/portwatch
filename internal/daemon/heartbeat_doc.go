// Package daemon — heartbeat
//
// heartbeat is a lightweight liveness tracker that the scan pipeline
// calls on every successful iteration. It records the wall-clock time
// of the most recent pulse and exposes helpers used by two consumers:
//
//  1. The health HTTP endpoint reads age() and count() to populate the
//     /healthz JSON response, giving operators a quick signal that the
//     daemon is actively scanning.
//
//  2. The watchdog can call alive(threshold) to decide whether to fire
//     its stale-loop callback — e.g. if no pulse arrives within 2×the
//     configured scan interval, the daemon is considered hung.
//
// All methods are safe for concurrent use.
package daemon
