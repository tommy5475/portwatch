// Package daemon — epoch
//
// epoch is a lightweight, concurrency-safe monotonic generation counter.
//
// # Purpose
//
// Many components in portwatch need to detect whether shared state has
// changed between two observations without retaining a full copy of that
// state. epoch provides a simple integer generation that is incremented
// every time a meaningful state transition occurs.
//
// # Usage
//
//	e := newEpoch()
//
//	// record the generation before doing work
//	gen := e.generation()
//
//	// … perform work …
//
//	// advance after a successful state change
//	e.advance()
//
//	// elsewhere: check whether state has changed
//	if e.generation() != gen {
//		// state changed; re-read it
//	}
//
// # Thread safety
//
// All methods are safe for concurrent use. Reads use a shared lock;
// advance uses an exclusive lock.
package daemon
