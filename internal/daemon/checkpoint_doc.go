// Package daemon provides internal daemon utilities for portwatch.
//
// # Checkpoint
//
// A checkpoint tracks named milestones reached during daemon operation.
// It is useful for observability — knowing which stages of the scan
// pipeline have been reached, how many times, and when.
//
// Example usage:
//
//	cp := newCheckpoint()
//	cp.Mark("scan.start")
//	cp.Mark("scan.complete")
//	// later…
//	fmt.Println(cp.Count("scan.complete"))  // 1
//	fmt.Println(cp.LastAt("scan.complete")) // time of last mark
//
// All methods are safe for concurrent use.
package daemon
