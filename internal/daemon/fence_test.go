package daemon

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestFenceInitiallyNotCrossed(t *testing.T) {
	f := newFence()
	if f.Crossed() {
		t.Fatal("expected fence to be uncrossed initially")
	}
}

func TestFenceCrossedAtZeroBeforeCross(t *testing.T) {
	f := newFence()
	if !f.CrossedAt().IsZero() {
		t.Fatal("expected CrossedAt to be zero before Cross")
	}
}

func TestFenceCrossSetsCrossed(t *testing.T) {
	f := newFence()
	f.Cross()
	if !f.Crossed() {
		t.Fatal("expected fence to be crossed after Cross()")
	}
}

func TestFenceCrossedAtIsSetAfterCross(t *testing.T) {
	before := time.Now()
	f := newFence()
	f.Cross()
	after := time.Now()
	at := f.CrossedAt()
	if at.Before(before) || at.After(after) {
		t.Fatalf("CrossedAt %v not in expected range [%v, %v]", at, before, after)
	}
}

func TestFenceCrossIsIdempotent(t *testing.T) {
	f := newFence()
	f.Cross()
	t1 := f.CrossedAt()
	time.Sleep(2 * time.Millisecond)
	f.Cross()
	t2 := f.CrossedAt()
	if !t1.Equal(t2) {
		t.Fatal("expected CrossedAt to be unchanged on second Cross")
	}
}

func TestFenceWaitReturnsTrueWhenCrossed(t *testing.T) {
	f := newFence()
	go func() {
		time.Sleep(10 * time.Millisecond)
		f.Cross()
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if !f.Wait(ctx) {
		t.Fatal("expected Wait to return true after Cross")
	}
}

func TestFenceWaitReturnsFalseOnContextCancel(t *testing.T) {
	f := newFence()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if f.Wait(ctx) {
		t.Fatal("expected Wait to return false when context already cancelled")
	}
}

func TestFenceWaitImmediatelyWhenAlreadyCrossed(t *testing.T) {
	f := newFence()
	f.Cross()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if !f.Wait(ctx) {
		t.Fatal("expected Wait to return true immediately for crossed fence")
	}
}

func TestFenceConcurrentWaiters(t *testing.T) {
	f := newFence()
	const n = 20
	var wg sync.WaitGroup
	results := make([]bool, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			results[idx] = f.Wait(ctx)
		}(i)
	}
	time.Sleep(5 * time.Millisecond)
	f.Cross()
	wg.Wait()
	for i, r := range results {
		if !r {
			t.Errorf("waiter %d did not receive cross signal", i)
		}
	}
}
