package daemon

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatcherInitialState(t *testing.T) {
	w := newWatcher(3, nil, nil)
	if w.isDegraded() {
		t.Fatal("expected watcher to start healthy")
	}
}

func TestWatcherDegradeOnThreshold(t *testing.T) {
	var degraded atomic.Bool
	w := newWatcher(3, func() { degraded.Store(true) }, nil)

	w.recordFailure()
	w.recordFailure()
	if degraded.Load() {
		t.Fatal("should not degrade before threshold")
	}
	w.recordFailure()
	if !degraded.Load() {
		t.Fatal("expected degrade callback after threshold")
	}
	if !w.isDegraded() {
		t.Fatal("expected isDegraded to return true")
	}
}

func TestWatcherRecoverAfterSuccess(t *testing.T) {
	var recovered atomic.Bool
	w := newWatcher(2, nil, func() { recovered.Store(true) })

	w.recordFailure()
	w.recordFailure()
	w.recordSuccess()

	if !recovered.Load() {
		t.Fatal("expected recover callback")
	}
	if w.isDegraded() {
		t.Fatal("expected healthy state after recovery")
	}
}

func TestWatcherSuccessWithoutPriorFailure(t *testing.T) {
	var recovered atomic.Bool
	w := newWatcher(3, nil, func() { recovered.Store(true) })
	w.recordSuccess()
	if recovered.Load() {
		t.Fatal("recover should not fire when never degraded")
	}
}

func TestWatcherLastSuccessTime(t *testing.T) {
	w := newWatcher(3, nil, nil)
	if !w.lastSuccessTime().IsZero() {
		t.Fatal("expected zero time before any success")
	}
	before := time.Now()
	w.recordSuccess()
	after := time.Now()
	ts := w.lastSuccessTime()
	if ts.Before(before) || ts.After(after) {
		t.Fatalf("unexpected last success time: %v", ts)
	}
}

func TestWatcherDefaultThreshold(t *testing.T) {
	w := newWatcher(0, nil, nil) // invalid → defaults to 3
	if w.threshold != 3 {
		t.Fatalf("expected default threshold 3, got %d", w.threshold)
	}
}

func TestWatcherRunUntilCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	w := newWatcher(3, nil, nil)
	calls := atomic.Int64{}
	errFn := func() error {
		calls.Add(1)
		return errors.New("fail")
	}
	go w.runUntil(ctx, 10*time.Millisecond, errFn)
	time.Sleep(55 * time.Millisecond)
	cancel()
	time.Sleep(15 * time.Millisecond)
	if calls.Load() < 3 {
		t.Fatalf("expected at least 3 calls, got %d", calls.Load())
	}
}
