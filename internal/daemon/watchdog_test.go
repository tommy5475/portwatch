package daemon

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestWatchdogDoesNotFireWhenKicked(t *testing.T) {
	t.Parallel()
	var fired atomic.Int32
	w := newWatchdog(80*time.Millisecond, func() { fired.Add(1) })
	defer w.Stop()

	for i := 0; i < 5; i++ {
		time.Sleep(20 * time.Millisecond)
		w.Kick()
	}

	if fired.Load() != 0 {
		t.Fatalf("expected watchdog not to fire, but it fired %d time(s)", fired.Load())
	}
}

func TestWatchdogFiresAfterTimeout(t *testing.T) {
	t.Parallel()
	ch := make(chan struct{}, 1)
	w := newWatchdog(40*time.Millisecond, func() { ch <- struct{}{} })
	defer w.Stop()

	select {
	case <-ch:
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("watchdog did not fire within expected window")
	}

	if !w.Fired() {
		t.Fatal("Fired() should return true after watchdog fires")
	}
}

func TestWatchdogKickClearsFiredState(t *testing.T) {
	t.Parallel()
	ch := make(chan struct{}, 1)
	w := newWatchdog(30*time.Millisecond, func() {
		select {
		case ch <- struct{}{}:
		default:
		}
	})
	defer w.Stop()

	<-ch // wait for first fire
	w.Kick()

	if w.Fired() {
		t.Fatal("Fired() should be false immediately after Kick")
	}
}

func TestWatchdogStopPreventsCallback(t *testing.T) {
	t.Parallel()
	var fired atomic.Int32
	w := newWatchdog(50*time.Millisecond, func() { fired.Add(1) })
	w.Stop()

	time.Sleep(120 * time.Millisecond)
	if fired.Load() != 0 {
		t.Fatalf("watchdog fired after Stop(); count=%d", fired.Load())
	}
}

func TestWatchdogDefaultsOnInvalidTimeout(t *testing.T) {
	t.Parallel()
	// Should not panic; zero timeout uses the default.
	w := newWatchdog(0, func() {})
	defer w.Stop()
	if w.timeout != 30*time.Second {
		t.Fatalf("expected default timeout 30s, got %v", w.timeout)
	}
}
