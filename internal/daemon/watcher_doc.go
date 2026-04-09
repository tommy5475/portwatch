// Package daemon provides the core runtime components for the portwatch
// daemon process.
//
// # Watcher
//
// The watcher subsystem observes the health of the periodic scan pipeline
// and transitions the daemon between healthy and degraded states based on
// consecutive failure counts.
//
// A watcher is created with a failure threshold and optional callback
// functions that are invoked exactly once when the system crosses into a
// degraded state, and again when it recovers:
//
//	w := newWatcher(3,
//		func() { log.Println("scan pipeline degraded") },
//		func() { log.Println("scan pipeline recovered") },
//	)
//
// The watcher integrates naturally with the health-check endpoint: the
// isDegraded method can be polled by the HTTP handler to surface a
// non-200 status code when the daemon is unhealthy.
//
// Thread safety: all watcher methods are safe for concurrent use.
package daemon
