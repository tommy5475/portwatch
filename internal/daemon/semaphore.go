package daemon

import "context"

// semaphore is a counting semaphore that limits concurrent access to a
// resource. It is safe for concurrent use by multiple goroutines.
type semaphore struct {
	ch chan struct{}
}

// newSemaphore creates a semaphore with the given concurrency limit.
// If n < 1, it defaults to 1.
func newSemaphore(n int) *semaphore {
	if n < 1 {
		n = 1
	}
	ch := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		ch <- struct{}{}
	}
	return &semaphore{ch: ch}
}

// Acquire blocks until a slot is available or ctx is cancelled.
// Returns an error if ctx is done before a slot is acquired.
func (s *semaphore) Acquire(ctx context.Context) error {
	select {
	case <-s.ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release returns a slot to the semaphore. It panics if called more
// times than Acquire has successfully returned.
func (s *semaphore) Release() {
	select {
	case s.ch <- struct{}{}:
	default:
		panic("semaphore: Release called without matching Acquire")
	}
}

// Available returns the number of slots currently available.
func (s *semaphore) Available() int {
	return len(s.ch)
}
