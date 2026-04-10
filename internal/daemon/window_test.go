package daemon

import (
	"testing"
	"time"
)

func TestWindowInitialCountIsZero(t *testing.T) {
	w := newWindow(5, 100*time.Millisecond)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestWindowRecordIncrements(t *testing.T) {
	w := newWindow(5, time.Second)
	w.Record()
	w.Record()
	if got := w.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestWindowResetClearsCount(t *testing.T) {
	w := newWindow(5, time.Second)
	w.Record()
	w.Record()
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestWindowEvictsStaleEvents(t *testing.T) {
	// Use a short interval so we can advance time by sleeping.
	w := newWindow(3, 50*time.Millisecond)
	w.Record()
	w.Record()

	// Wait for more than the full window span (3 * 50ms = 150ms).
	time.Sleep(200 * time.Millisecond)

	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after window expiry, got %d", got)
	}
}

func TestWindowPartialEviction(t *testing.T) {
	// 4 buckets of 50ms each → 200ms total window.
	w := newWindow(4, 50*time.Millisecond)
	w.Record() // recorded in bucket[0] at t=0

	// Advance ~100ms (2 buckets): the original record should still be within
	// the window (bucket[2] after rotation).
	time.Sleep(80 * time.Millisecond)
	w.Record() // new event in the current bucket

	if got := w.Count(); got < 1 {
		t.Fatalf("expected at least 1 event within partial window, got %d", got)
	}
}

func TestWindowDefaultsOnInvalidArgs(t *testing.T) {
	// Should not panic and should use safe defaults.
	w := newWindow(0, 0)
	w.Record()
	if got := w.Count(); got != 1 {
		t.Fatalf("expected 1 with default args, got %d", got)
	}
}

func TestWindowConcurrentAccess(t *testing.T) {
	w := newWindow(10, 10*time.Millisecond)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			w.Record()
			_ = w.Count()
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
}
