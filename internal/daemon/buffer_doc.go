// Package daemon contains the core runtime components of portwatch.
//
// # Buffer
//
// buffer is a generic, bounded ring buffer safe for concurrent use.
// It retains the most recent N items and silently overwrites the oldest
// entry once the capacity limit is reached.
//
// Typical usage — keep a rolling window of the last 100 scan events:
//
//		b := newBuffer[ScanEvent](100)
//		b.Push(event)
//		events := b.Snapshot() // ordered oldest → newest
//
// The zero value is not usable; always construct via newBuffer.
package daemon
