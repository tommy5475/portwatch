// Package daemon contains the runtime components of portwatch.
//
// # Semaphore
//
// semaphore is a lightweight counting semaphore used to bound the
// number of goroutines that may concurrently execute a critical
// section — most notably parallel port-scan workers.
//
// Usage:
//
//	sem := newSemaphore(4) // allow up to 4 concurrent workers
//
//	for _, port := range ports {
//		port := port
//		go func() {
//			if err := sem.Acquire(ctx); err != nil {
//				return // context cancelled
//			}
//			defer sem.Release()
//			scanPort(port)
//		}()
//	}
//
// The semaphore is context-aware: if the supplied context is cancelled
// while a goroutine is waiting for a slot, Acquire returns immediately
// with the context error so the caller can propagate cancellation
// without leaking goroutines.
package daemon
