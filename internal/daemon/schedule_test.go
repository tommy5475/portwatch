package daemon

import (
	"testing"
	"time"
)

func TestScheduleDefaultsOnInvalidArgs(t *testing.T) {
	_, err := newSchedule(-1*time.Second, 10*time.Second, 1.5)
	if err == nil {
		t.Fatal("expected error for non-positive nominal")
	}
}

func TestScheduleInitialIntervalIsNominal(t *testing.T) {
	s, _ := newSchedule(5*time.Second, 60*time.Second, 2.0)
	if s.interval() != 5*time.Second {
		t.Fatalf("expected 5s, got %v", s.interval())
	}
}

func TestScheduleNextSuccessKeepsNominal(t *testing.T) {
	s, _ := newSchedule(5*time.Second, 60*time.Second, 2.0)
	d := s.next(false)
	if d != 5*time.Second {
		t.Fatalf("expected 5s, got %v", d)
	}
	if s.interval() != 5*time.Second {
		t.Fatalf("interval should remain nominal after success")
	}
}

func TestScheduleNextFailureDoublesInterval(t *testing.T) {
	s, _ := newSchedule(5*time.Second, 60*time.Second, 2.0)
	s.next(true)
	if s.interval() != 10*time.Second {
		t.Fatalf("expected 10s after one failure, got %v", s.interval())
	}
}

func TestScheduleCapsAtMaxDelay(t *testing.T) {
	s, _ := newSchedule(5*time.Second, 12*time.Second, 2.0)
	s.next(true) // 10s
	s.next(true) // would be 20s, capped at 12s
	if s.interval() != 12*time.Second {
		t.Fatalf("expected cap at 12s, got %v", s.interval())
	}
}

func TestScheduleResetRestoresNominal(t *testing.T) {
	s, _ := newSchedule(5*time.Second, 60*time.Second, 2.0)
	s.next(true)
	s.next(true)
	s.reset()
	if s.interval() != 5*time.Second {
		t.Fatalf("expected nominal 5s after reset, got %v", s.interval())
	}
}

func TestScheduleMaxDelayClampedToNominal(t *testing.T) {
	s, _ := newSchedule(10*time.Second, 1*time.Second, 2.0)
	// maxDelay < nominal should be clamped to nominal
	if s.maxDelay != 10*time.Second {
		t.Fatalf("expected maxDelay clamped to 10s, got %v", s.maxDelay)
	}
}

func TestScheduleInvalidFactorDefaulted(t *testing.T) {
	s, _ := newSchedule(5*time.Second, 60*time.Second, 0.5)
	s.next(true)
	// factor defaulted to 1.5 => 5s * 1.5 = 7.5s
	if s.interval() != time.Duration(float64(5*time.Second)*1.5) {
		t.Fatalf("unexpected interval after invalid factor default: %v", s.interval())
	}
}
