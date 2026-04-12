// Package daemon provides the core runtime for portwatch.
//
// # Hooks
//
// The hooks subsystem provides a lightweight lifecycle-event mechanism that
// allows internal components (and tests) to react to daemon state transitions
// without introducing circular dependencies.
//
// Supported events:
//
//	- HookBeforeScan   – fired immediately before each port-scan cycle.
//	- HookAfterScan    – fired after a scan completes, payload is scan duration.
//	- HookOnChange     – fired when the port diff is non-empty.
//	- HookOnDegraded   – fired when the watcher enters degraded mode.
//	- HookOnRecovered  – fired when the watcher recovers from degraded mode.
//
// Usage:
//
//	 h := newHooks()
//	 h.Register(HookOnChange, func(e hookEvent, p any) {
//	     log.Printf("ports changed: %v", p)
//	 })
//	 h.Fire(HookOnChange, diff)
package daemon
