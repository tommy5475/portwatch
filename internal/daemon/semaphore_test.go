package daemon

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSemaphoreDefaultsToOne(t *testing.T) {
	sem := newSemaphore(0)
	if sem.Available() != 1 {
		t.Fatalf("expected 1 slot, got %d", sem.Available())
	}
}

func TestSemaphoreAcquireRelease(t *testing.T) {
	sem := newSemaphore(2)
	ctx := context.Background()

	if err := sem.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sem.Available() != 1 {
		t.Fatalf("expected 1 slot after acquire, got %d", sem.Available())
	}
	sem.Release()
	if sem.Available() != 2 {
		t.Fatalf("expected 2 slots after release, got %d", sem.Available())
	}
}

func TestSemaphoreBlocksWhenFull(t *testing.T) {
	sem := newSemaphore(1)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Drain the single slot.
	_ = sem.Acquire(context.Background())

	// Second acquire must block and eventually fail due to timeout.
	if err := sem.Acquire(ctx); err == nil {
		t.Fatal("expected error when semaphore is full, got nil")
	}
}

func TestSemaphoreConcurrentWorkers(t *testing.T) {
	const limit = 3
	const workers = 12

	sem := newSemaphore(limit)
	ctx := context.Background()

	var active atomic.Int32
	var maxSeen atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sem.Acquire(ctx)
			defer sem.Release()

			cur := active.Add(1)
			for {
				old := maxSeen.Load()
				if cur <= old || maxSeen.CompareAndSwap(old, cur) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			active.Add(-1)
		}()
	}
	wg.Wait()

	if got := maxSeen.Load(); got > limit {
		t.Fatalf("concurrency limit exceeded: max active = %d, limit = %d", got, limit)
	}
}

func TestSemaphoreReleasePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on over-release")
		}
	}()
	sem := newSemaphore(1)
	sem.Release() // one too many — should panic
}
