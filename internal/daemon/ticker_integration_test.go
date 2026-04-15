package daemon

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunLoopIntegrationTicksMultipleTimes(t *testing.T) {
	t.Parallel()

	var count atomic.Int64
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	ticker := newTicker(50 * time.Millisecond)
	defer ticker.stop()

	done := make(chan struct{})
	go func() {
		defer close(done)
		runLoop(ctx, ticker, func() {
			count.Add(1)
		})
	}()

	<-done

	got := count.Load()
	if got < 3 {
		t.Errorf("expected at least 3 ticks in 300ms with 50ms interval, got %d", got)
	}
}

func TestRunLoopIntegrationStopsCleanlyOnCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	ticker := newTicker(10 * time.Millisecond)
	defer ticker.stop()

	done := make(chan struct{})
	go func() {
		defer close(done)
		runLoop(ctx, ticker, func() {})
	}()

	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("runLoop did not stop after context cancellation")
	}
}

func TestRunLoopIntegrationFnBlockDoesNotMissCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	ticker := newTicker(10 * time.Millisecond)
	defer ticker.stop()

	started := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		runLoop(ctx, ticker, func() {
			select {
			case <-started:
			default:
				close(started)
			}
			time.Sleep(20 * time.Millisecond)
		})
	}()

	<-started
	cancel()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("runLoop did not exit after fn completed and context was cancelled")
	}
}
