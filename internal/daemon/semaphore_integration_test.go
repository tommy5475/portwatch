package daemon

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestSemaphoreStormUnderLoad fires a large number of goroutines and
// verifies that the semaphore correctly gates them to the configured
// concurrency limit under sustained load.
func TestSemaphoreStormUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	const limit = 5
	const workers = 100

	sem := newSemaphore(limit)
	ctx := context.Background()

	var active atomic.Int32
	var violations atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := sem.Acquire(ctx); err != nil {
				return
			}
			defer sem.Release()

			if active.Add(1) > limit {
				violations.Add(1)
			}
			time.Sleep(2 * time.Millisecond)
			active.Add(-1)
		}()
	}
	wg.Wait()

	if v := violations.Load(); v > 0 {
		t.Fatalf("detected %d concurrency-limit violations", v)
	}
}

// TestSemaphoreContextCancelDrainsWaiters ensures that cancelling a
// context unblocks all goroutines waiting on Acquire.
func TestSemaphoreContextCancelDrainsWaiters(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	sem := newSemaphore(1)
	// Hold the only slot.
	_ = sem.Acquire(context.Background())

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sem.Acquire(ctx) // will block until cancel
		}()
	}

	time.Sleep(20 * time.Millisecond)
	cancel()

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
		// all waiters unblocked — success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for goroutines to unblock after cancel")
	}
}
