package daemon

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// TestWatcherIntegrationDegradeAndRecover exercises the full degrade →
// recover lifecycle using runUntil with a real ticker.
func TestWatcherIntegrationDegradeAndRecover(t *testing.T) {
	var (
		degradeCount atomic.Int64
		recoverCount atomic.Int64
	)
	w := newWatcher(3,
		func() { degradeCount.Add(1) },
		func() { recoverCount.Add(1) },
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Phase 1: always fail → should degrade after 3 ticks.
	failPhase := atomic.Bool{}
	failPhase.Store(true)

	checkFn := func() error {
		if failPhase.Load() {
			return errors.New("simulated error")
		}
		return nil
	}

	go w.runUntil(ctx, 20*time.Millisecond, checkFn)

	// Wait long enough for at least 3 failures.
	time.Sleep(100 * time.Millisecond)

	if !w.isDegraded() {
		t.Fatal("expected degraded state after repeated failures")
	}
	if degradeCount.Load() != 1 {
		t.Fatalf("expected exactly 1 degrade event, got %d", degradeCount.Load())
	}

	// Phase 2: switch to success → should recover.
	failPhase.Store(false)
	time.Sleep(60 * time.Millisecond)

	if w.isDegraded() {
		t.Fatal("expected healthy state after recovery")
	}
	if recoverCount.Load() != 1 {
		t.Fatalf("expected exactly 1 recover event, got %d", recoverCount.Load())
	}
}
