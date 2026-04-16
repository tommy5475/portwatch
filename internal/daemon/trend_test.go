package daemon

import (
	"testing"
	"time"
)

func TestTrendInitialRateIsZero(t *testing.T) {
	tr := newTrend(10, time.Minute)
	if tr.rate() != 0 {
		t.Fatalf("expected 0, got %d", tr.rate())
	}
}

func TestTrendRecordIncreasesRate(t *testing.T) {
	tr := newTrend(10, time.Minute)
	tr.record(5)
	tr.record(3)
	if tr.rate() < 8 {
		t.Fatalf("expected at least 8, got %d", tr.rate())
	}
}

func TestTrendResetClearsRate(t *testing.T) {
	tr := newTrend(10, time.Minute)
	tr.record(10)
	tr.reset()
	if tr.rate() != 0 {
		t.Fatalf("expected 0 after reset, got %d", tr.rate())
	}
}

func TestTrendRisingWhenSecondHalfHigher(t *testing.T) {
	tr := newTrend(4, time.Second)
	// fill first half (buckets 0,1) with low values via direct write
	tr.mu.Lock()
	tr.buckets[0] = 1
	tr.buckets[1] = 1
	tr.buckets[2] = 10
	tr.buckets[3] = 10
	tr.mu.Unlock()
	if !tr.rising() {
		t.Fatal("expected rising to be true")
	}
}

func TestTrendNotRisingWhenFirstHalfHigher(t *testing.T) {
	tr := newTrend(4, time.Second)
	tr.mu.Lock()
	tr.buckets[0] = 10
	tr.buckets[1] = 10
	tr.buckets[2] = 1
	tr.buckets[3] = 1
	tr.mu.Unlock()
	if tr.rising() {
		t.Fatal("expected rising to be false")
	}
}

func TestTrendDefaultsOnInvalidArgs(t *testing.T) {
	tr := newTrend(0, -1)
	if tr.size < 2 {
		t.Fatalf("expected default size >= 2, got %d", tr.size)
	}
	if tr.interval <= 0 {
		t.Fatal("expected positive default interval")
	}
}

func TestTrendBucketIndexInBounds(t *testing.T) {
	tr := newTrend(8, time.Second)
	for i := 0; i < 20; i++ {
		idx := tr.bucketIndex(time.Now().Add(time.Duration(i) * time.Second))
		if idx < 0 || idx >= tr.size {
			t.Fatalf("bucket index %d out of bounds [0,%d)", idx, tr.size)
		}
	}
}
