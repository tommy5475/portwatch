package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestLimiterAllowsUpToMax(t *testing.T) {
	l := newLimiter(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !l.Allow("k") {
			t.Fatalf("expected Allow to return true on call %d", i+1)
		}
	}
	if l.Allow("k") {
		t.Fatal("expected Allow to return false after max reached")
	}
}

func TestLimiterRemainingDecreases(t *testing.T) {
	l := newLimiter(5, time.Minute)
	if r := l.Remaining("k"); r != 5 {
		t.Fatalf("expected 5 remaining, got %d", r)
	}
	l.Allow("k")
	l.Allow("k")
	if r := l.Remaining("k"); r != 3 {
		t.Fatalf("expected 3 remaining, got %d", r)
	}
}

func TestLimiterEvictsStaleEntries(t *testing.T) {
	now := time.Now()
	l := newLimiter(2, time.Second)
	// Inject a timestamp older than the window.
	l.buckets["k"] = []time.Time{now.Add(-2 * time.Second)}
	if !l.Allow("k") {
		t.Fatal("expected stale entry to be evicted and Allow to succeed")
	}
}

func TestLimiterReset(t *testing.T) {
	l := newLimiter(2, time.Minute)
	l.Allow("k")
	l.Allow("k")
	if l.Allow("k") {
		t.Fatal("expected limit to be reached before reset")
	}
	l.Reset("k")
	if !l.Allow("k") {
		t.Fatal("expected Allow to succeed after reset")
	}
}

func TestLimiterKeyIsolation(t *testing.T) {
	l := newLimiter(1, time.Minute)
	l.Allow("a")
	if !l.Allow("b") {
		t.Fatal("key 'b' should be independent of key 'a'")
	}
}

func TestLimiterDefaultsOnInvalidArgs(t *testing.T) {
	l := newLimiter(0, 0)
	if l.max != 10 {
		t.Fatalf("expected default max 10, got %d", l.max)
	}
	if l.window != time.Minute {
		t.Fatalf("expected default window 1m, got %v", l.window)
	}
}

func TestLimiterConcurrentAccess(t *testing.T) {
	l := newLimiter(50, time.Minute)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Allow("shared")
		}()
	}
	wg.Wait()
	if r := l.Remaining("shared"); r != 0 {
		t.Fatalf("expected 0 remaining after 100 concurrent calls with max 50, got %d", r)
	}
}
