package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestDialDefaultTarget(t *testing.T) {
	d := newDial("")
	if d.target != "unknown" {
		t.Fatalf("expected 'unknown', got %q", d.target)
	}
}

func TestDialInitialCountIsZero(t *testing.T) {
	d := newDial("localhost:9000")
	if d.count() != 0 {
		t.Fatalf("expected 0, got %d", d.count())
	}
}

func TestDialSuccessRateNoAttempts(t *testing.T) {
	d := newDial("localhost:9000")
	if d.successRate() != 0 {
		t.Fatalf("expected 0.0, got %f", d.successRate())
	}
}

func TestDialRecordSuccess(t *testing.T) {
	d := newDial("localhost:9000")
	d.record(true, 5*time.Millisecond)
	if d.count() != 1 {
		t.Fatalf("expected 1, got %d", d.count())
	}
	if d.successRate() != 1.0 {
		t.Fatalf("expected 1.0, got %f", d.successRate())
	}
	if !d.lastSuccess() {
		t.Fatal("expected last attempt to be success")
	}
}

func TestDialRecordFailure(t *testing.T) {
	d := newDial("localhost:9000")
	d.record(false, 2*time.Millisecond)
	if d.successRate() != 0 {
		t.Fatalf("expected 0.0, got %f", d.successRate())
	}
	if d.lastSuccess() {
		t.Fatal("expected last attempt to be failure")
	}
}

func TestDialAvgLatency(t *testing.T) {
	d := newDial("localhost:9000")
	d.record(true, 10*time.Millisecond)
	d.record(true, 20*time.Millisecond)
	got := d.avgLatency()
	if got != 15*time.Millisecond {
		t.Fatalf("expected 15ms, got %v", got)
	}
}

func TestDialNegativeLatencyClamped(t *testing.T) {
	d := newDial("localhost:9000")
	d.record(true, -5*time.Millisecond)
	if d.avgLatency() != 0 {
		t.Fatalf("expected 0, got %v", d.avgLatency())
	}
}

func TestDialReset(t *testing.T) {
	d := newDial("localhost:9000")
	d.record(true, 10*time.Millisecond)
	d.reset()
	if d.count() != 0 {
		t.Fatalf("expected 0 after reset, got %d", d.count())
	}
	if d.successRate() != 0 {
		t.Fatalf("expected 0.0 after reset")
	}
}

func TestDialConcurrentRecord(t *testing.T) {
	d := newDial("localhost:9000")
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			d.record(i%2 == 0, time.Millisecond)
		}(i)
	}
	wg.Wait()
	if d.count() != 50 {
		t.Fatalf("expected 50, got %d", d.count())
	}
}
