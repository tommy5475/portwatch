package daemon

import (
	"testing"
	"time"
)

func TestTickerMetricsInitialState(t *testing.T) {
	m := newTickerMetrics()
	if m.TickCount() != 0 {
		t.Fatalf("expected 0 tick count, got %d", m.TickCount())
	}
	if m.SkipCount() != 0 {
		t.Fatalf("expected 0 skip count, got %d", m.SkipCount())
	}
	if !m.LastTickedAt().IsZero() {
		t.Fatalf("expected zero LastTickedAt")
	}
	if m.AverageLag() != 0 {
		t.Fatalf("expected zero average lag")
	}
}

func TestTickerMetricsRecordTick(t *testing.T) {
	m := newTickerMetrics()
	now := time.Now()
	lag := 5 * time.Millisecond
	m.recordTick(now, lag)

	if m.TickCount() != 1 {
		t.Fatalf("expected tick count 1, got %d", m.TickCount())
	}
	if m.LastTickedAt().IsZero() {
		t.Fatal("expected non-zero LastTickedAt")
	}
	if m.AverageLag() != lag {
		t.Fatalf("expected average lag %v, got %v", lag, m.AverageLag())
	}
}

func TestTickerMetricsRecordSkip(t *testing.T) {
	m := newTickerMetrics()
	m.recordSkip()
	m.recordSkip()
	if m.SkipCount() != 2 {
		t.Fatalf("expected skip count 2, got %d", m.SkipCount())
	}
}

func TestTickerMetricsAverageLagMultipleTicks(t *testing.T) {
	m := newTickerMetrics()
	m.recordTick(time.Now(), 10*time.Millisecond)
	m.recordTick(time.Now(), 20*time.Millisecond)
	expected := 15 * time.Millisecond
	if m.AverageLag() != expected {
		t.Fatalf("expected average lag %v, got %v", expected, m.AverageLag())
	}
}

func TestTickerMetricsNegativeLagIgnored(t *testing.T) {
	m := newTickerMetrics()
	m.recordTick(time.Now(), -1*time.Millisecond)
	if m.totalLag.Load() != 0 {
		t.Fatal("negative lag should not be recorded")
	}
}

func TestTickerMetricsReset(t *testing.T) {
	m := newTickerMetrics()
	m.recordTick(time.Now(), 5*time.Millisecond)
	m.recordSkip()
	m.reset()
	if m.TickCount() != 0 || m.SkipCount() != 0 || !m.LastTickedAt().IsZero() {
		t.Fatal("reset did not clear metrics")
	}
}
