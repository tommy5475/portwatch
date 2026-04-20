package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestBenchInitialState(t *testing.T) {
	b := newBench()
	if b.len() != 0 {
		t.Fatalf("expected count 0, got %d", b.len())
	}
	if b.mean() != 0 {
		t.Fatalf("expected mean 0, got %v", b.mean())
	}
	if b.minDuration() != 0 {
		t.Fatalf("expected min 0, got %v", b.minDuration())
	}
	if b.maxDuration() != 0 {
		t.Fatalf("expected max 0, got %v", b.maxDuration())
	}
}

func TestBenchRecordSingle(t *testing.T) {
	b := newBench()
	b.record(10 * time.Millisecond)
	if b.len() != 1 {
		t.Fatalf("expected count 1, got %d", b.len())
	}
	if b.mean() != 10*time.Millisecond {
		t.Fatalf("expected mean 10ms, got %v", b.mean())
	}
	if b.minDuration() != 10*time.Millisecond {
		t.Fatalf("expected min 10ms, got %v", b.minDuration())
	}
	if b.maxDuration() != 10*time.Millisecond {
		t.Fatalf("expected max 10ms, got %v", b.maxDuration())
	}
}

func TestBenchRecordMultiple(t *testing.T) {
	b := newBench()
	b.record(10 * time.Millisecond)
	b.record(30 * time.Millisecond)
	b.record(20 * time.Millisecond)
	if b.len() != 3 {
		t.Fatalf("expected count 3, got %d", b.len())
	}
	if b.mean() != 20*time.Millisecond {
		t.Fatalf("expected mean 20ms, got %v", b.mean())
	}
	if b.minDuration() != 10*time.Millisecond {
		t.Fatalf("expected min 10ms, got %v", b.minDuration())
	}
	if b.maxDuration() != 30*time.Millisecond {
		t.Fatalf("expected max 30ms, got %v", b.maxDuration())
	}
}

func TestBenchNegativeDurationClampedToZero(t *testing.T) {
	b := newBench()
	b.record(-5 * time.Millisecond)
	if b.minDuration() != 0 {
		t.Fatalf("expected min 0 for negative input, got %v", b.minDuration())
	}
}

func TestBenchReset(t *testing.T) {
	b := newBench()
	b.record(10 * time.Millisecond)
	b.record(20 * time.Millisecond)
	b.reset()
	if b.len() != 0 {
		t.Fatalf("expected count 0 after reset, got %d", b.len())
	}
	if b.mean() != 0 {
		t.Fatalf("expected mean 0 after reset, got %v", b.mean())
	}
}

func TestBenchConcurrentRecord(t *testing.T) {
	b := newBench()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.record(time.Millisecond)
		}()
	}
	wg.Wait()
	if b.len() != 100 {
		t.Fatalf("expected count 100, got %d", b.len())
	}
}
