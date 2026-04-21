package daemon

import (
	"testing"
	"time"
)

func TestSamplerInitialState(t *testing.T) {
	s := newSampler()
	st := s.stats()
	if st.Count != 0 {
		t.Fatalf("expected count 0, got %d", st.Count)
	}
	if st.Mean != 0 {
		t.Fatalf("expected mean 0, got %f", st.Mean)
	}
}

func TestSamplerRecordUpdatesCount(t *testing.T) {
	s := newSampler()
	s.record(10)
	s.record(20)
	if s.stats().Count != 2 {
		t.Fatalf("expected count 2")
	}
}

func TestSamplerMean(t *testing.T) {
	s := newSampler()
	s.record(10)
	s.record(20)
	s.record(30)
	got := s.mean()
	if got != 20 {
		t.Fatalf("expected mean 20, got %f", got)
	}
}

func TestSamplerMinMax(t *testing.T) {
	s := newSampler()
	s.record(5)
	s.record(15)
	s.record(10)
	st := s.stats()
	if st.Min != 5 {
		t.Fatalf("expected min 5, got %f", st.Min)
	}
	if st.Max != 15 {
		t.Fatalf("expected max 15, got %f", st.Max)
	}
}

func TestSamplerLastValue(t *testing.T) {
	s := newSampler()
	s.record(1)
	s.record(99)
	if s.stats().Last != 99 {
		t.Fatalf("expected last 99")
	}
}

func TestSamplerLastAtUpdated(t *testing.T) {
	s := newSampler()
	before := time.Now()
	s.record(7)
	after := time.Now()
	st := s.stats()
	if st.LastAt.Before(before) || st.LastAt.After(after) {
		t.Fatalf("LastAt out of expected range")
	}
}

func TestSamplerReset(t *testing.T) {
	s := newSampler()
	s.record(42)
	s.reset()
	st := s.stats()
	if st.Count != 0 {
		t.Fatalf("expected count 0 after reset")
	}
	if st.Mean != 0 {
		t.Fatalf("expected mean 0 after reset")
	}
}

func TestSamplerUptimePositive(t *testing.T) {
	s := newSampler()
	time.Sleep(time.Millisecond)
	if s.stats().Uptime <= 0 {
		t.Fatalf("expected positive uptime")
	}
}
