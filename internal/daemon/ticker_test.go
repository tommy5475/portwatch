package daemon

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

// fakeTicker is a controllable tickSource for unit tests.
type fakeTicker struct {
	ch   chan time.Time
	stopped atomic.Bool
}

func newFakeTicker() *fakeTicker {
	return &fakeTicker{ch: make(chan time.Time, 4)}
}

func (f *fakeTicker) C() <-chan time.Time { return f.ch }
func (f *fakeTicker) Stop()              { f.stopped.Store(true) }
func (f *fakeTicker) tick()              { f.ch <- time.Now() }

func TestRunLoopCallsFnOnTick(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ft := newFakeTicker()
	var count atomic.Int32

	done := make(chan error, 1)
	go func() {
		done <- runLoop(ctx, ft, func() { count.Add(1) })
	}()

	ft.tick()
	ft.tick()
	ft.tick()

	// Give the goroutine time to process all ticks.
	time.Sleep(20 * time.Millisecond)
	cancel()

	if err := <-done; err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if got := count.Load(); got != 3 {
		t.Fatalf("expected fn called 3 times, got %d", got)
	}
	if !ft.stopped.Load() {
		t.Fatal("expected ticker to be stopped after runLoop returns")
	}
}

func TestRunLoopStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ft := newFakeTicker()

	done := make(chan error, 1)
	go func() {
		done <- runLoop(ctx, ft, func() {})
	}()

	cancel()

	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("runLoop did not stop after context cancellation")
	}
}

func TestNewTicker(t *testing.T) {
	tk := newTicker(50 * time.Millisecond)
	defer tk.Stop()

	select {
	case <-tk.C():
		// received a tick — ticker is working
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected a tick within 200ms")
	}
}
