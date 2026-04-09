package daemon

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

var errFake = errors.New("fake error")

func TestRetrySucceedsOnFirstAttempt(t *testing.T) {
	p := newRetryPolicy(3, 10*time.Millisecond, 100*time.Millisecond)
	res := p.run(context.Background(), func() error { return nil })
	if res.err != nil {
		t.Fatalf("expected no error, got %v", res.err)
	}
	if res.attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", res.attempts)
	}
}

func TestRetryExhaustsAllAttempts(t *testing.T) {
	p := newRetryPolicy(3, 5*time.Millisecond, 20*time.Millisecond)
	var calls int32
	res := p.run(context.Background(), func() error {
		atomic.AddInt32(&calls, 1)
		return errFake
	})
	if res.err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
	if res.attempts != 3 {
		t.Fatalf("expected attempts=3, got %d", res.attempts)
	}
}

func TestRetrySucceedsOnSecondAttempt(t *testing.T) {
	p := newRetryPolicy(3, 5*time.Millisecond, 20*time.Millisecond)
	var calls int32
	res := p.run(context.Background(), func() error {
		if atomic.AddInt32(&calls, 1) < 2 {
			return errFake
		}
		return nil
	})
	if res.err != nil {
		t.Fatalf("expected success, got %v", res.err)
	}
	if res.attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", res.attempts)
	}
}

func TestRetryRespectsContextCancellation(t *testing.T) {
	p := newRetryPolicy(10, 50*time.Millisecond, 500*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	var calls int32
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	res := p.run(ctx, func() error {
		atomic.AddInt32(&calls, 1)
		return errFake
	})

	if res.err == nil {
		t.Fatal("expected context cancellation error")
	}
	if calls >= 10 {
		t.Fatalf("expected context to cut short retries, got %d calls", calls)
	}
}

func TestRetryDefaultsOnInvalidArgs(t *testing.T) {
	p := newRetryPolicy(0, -1, -1)
	if p.maxAttempts != 1 {
		t.Errorf("expected maxAttempts=1, got %d", p.maxAttempts)
	}
	if p.baseDelay <= 0 {
		t.Errorf("expected positive baseDelay, got %v", p.baseDelay)
	}
	if p.maxDelay < p.baseDelay {
		t.Errorf("maxDelay %v should be >= baseDelay %v", p.maxDelay, p.baseDelay)
	}
}
