// Package daemon provides internal runtime primitives used by the portwatch
// daemon, including scheduling, resilience, and coordination utilities.
//
// # Registry
//
// registry is a thread-safe key-value store for named runtime components.
// Each entry may carry an optional set of string tags that allow callers
// to query subsets of registered values by tag membership.
//
// Typical usage:
//
//	r := newRegistry()
//	r.Register("scanner", myScanner, "core", "active")
//	r.Register("alerter", myAlerter, "core")
//
//	coreNames := r.FilterByTag("core") // ["scanner", "alerter"]
//	val, err  := r.Get("scanner")
//
// The registry does not manage lifecycle; callers are responsible for
// starting and stopping any registered components.
package daemon
