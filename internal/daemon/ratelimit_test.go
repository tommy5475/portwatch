package daemon

import (
	"testing"
	"time"
)

func TestRateLimiterAllowsUpToMax(t *testing.T) {
	rl := newRateLimiter(3, time.Minute)

	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Fatalf("expected Allow() == true on call %d", i+1)
		}
	}
	if rl.Allow() {
		t.Fatal("expected Allow() == false after exhausting tokens")
	}
}

func TestRateLimiterRemainingDecreases(t *testing.T) {
	rl := newRateLimiter(5, time.Minute)
	if rl.Remaining() != 5 {
		t.Fatalf("expected 5 tokens, got %d", rl.Remaining())
	}
	rl.Allow()
	if rl.Remaining() != 4 {
		t.Fatalf("expected 4 tokens after one Allow(), got %d", rl.Remaining())
	}
}

func TestRateLimiterReplenishesAfterInterval(t *testing.T) {
	rl := newRateLimiter(2, 50*time.Millisecond)
	rl.Allow()
	rl.Allow()
	if rl.Allow() {
		t.Fatal("expected tokens to be exhausted")
	}

	time.Sleep(60 * time.Millisecond)

	if !rl.Allow() {
		t.Fatal("expected tokens to be replenished after interval")
	}
}

func TestRateLimiterDefaultsOnInvalidArgs(t *testing.T) {
	rl := newRateLimiter(0, 0)
	if rl.max != 10 {
		t.Fatalf("expected default max=10, got %d", rl.max)
	}
	if rl.interval != time.Minute {
		t.Fatalf("expected default interval=1m, got %v", rl.interval)
	}
}

func TestRateLimiterConcurrentAccess(t *testing.T) {
	rl := newRateLimiter(100, time.Minute)
	done := make(chan struct{})
	allowed := make(chan bool, 200)

	for i := 0; i < 200; i++ {
		go func() {
			allowed <- rl.Allow()
			done <- struct{}{}
		}()
	}
	for i := 0; i < 200; i++ {
		<-done
	}
	close(allowed)

	count := 0
	for v := range allowed {
		if v {
			count++
		}
	}
	if count != 100 {
		t.Fatalf("expected exactly 100 allowed events, got %d", count)
	}
}
