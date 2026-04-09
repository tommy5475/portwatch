package daemon

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestDebounceCallsFnAfterQuietPeriod(t *testing.T) {
	var count atomic.Int32
	d := newDebounce(50*time.Millisecond, func() { count.Add(1) })

	d.Trigger()
	time.Sleep(100 * time.Millisecond)

	if got := count.Load(); got != 1 {
		t.Fatalf("expected fn called once, got %d", got)
	}
}

func TestDebounceResetOnRapidTriggers(t *testing.T) {
	var count atomic.Int32
	d := newDebounce(60*time.Millisecond, func() { count.Add(1) })

	// Fire three times in quick succession — fn should run only once.
	for i := 0; i < 3; i++ {
		d.Trigger()
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(120 * time.Millisecond)

	if got := count.Load(); got != 1 {
		t.Fatalf("expected fn called once after burst, got %d", got)
	}
}

func TestDebounceFlushCallsImmediately(t *testing.T) {
	var count atomic.Int32
	d := newDebounce(500*time.Millisecond, func() { count.Add(1) })

	d.Trigger()
	d.Flush()
	time.Sleep(20 * time.Millisecond) // let goroutine run

	if got := count.Load(); got != 1 {
		t.Fatalf("expected fn called once via Flush, got %d", got)
	}
}

func TestDebounceStopPreventsCall(t *testing.T) {
	var count atomic.Int32
	d := newDebounce(50*time.Millisecond, func() { count.Add(1) })

	d.Trigger()
	d.Stop()
	time.Sleep(100 * time.Millisecond)

	if got := count.Load(); got != 0 {
		t.Fatalf("expected fn not called after Stop, got %d", got)
	}
}

func TestDebounceDefaultsOnInvalidWait(t *testing.T) {
	var count atomic.Int32
	d := newDebounce(-1, func() { count.Add(1) })

	if d.wait != 500*time.Millisecond {
		t.Fatalf("expected default wait 500ms, got %v", d.wait)
	}
	d.Stop()
	_ = count.Load()
}

func TestDebounceMultipleFlushIdempotent(t *testing.T) {
	var count atomic.Int32
	d := newDebounce(500*time.Millisecond, func() { count.Add(1) })

	d.Trigger()
	d.Flush()
	d.Flush() // second flush should be a no-op
	time.Sleep(20 * time.Millisecond)

	if got := count.Load(); got != 1 {
		t.Fatalf("expected fn called once, got %d", got)
	}
}
