// Package daemon provides the internal runtime primitives for the portwatch
// daemon process.
//
// # Trend
//
// trend is a lightweight rolling-rate tracker that partitions event counts
// into fixed-width time buckets. It is useful for detecting whether port-change
// activity is accelerating or settling down over a configurable observation
// window.
//
// Usage:
//
//	tr := newTrend(10, time.Minute)
//	tr.record(1)          // called each time a change is detected
//	if tr.rising() { … } // true when recent half outpaces earlier half
//	rate := tr.rate()     // total events across all buckets
//	tr.reset()            // clear all buckets (e.g. after a report)
//
// trend is safe for concurrent use.
package daemon
