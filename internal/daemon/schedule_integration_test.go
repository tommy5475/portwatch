package daemon

import (
	"testing"
	"time"
)

// TestScheduleFullBackoffAndRecoveryCycle exercises a realistic
// failure-then-recovery sequence end-to-end.
func TestScheduleFullBackoffAndRecoveryCycle(t *testing.T) {
	nominal := 5 * time.Second
	max := 40 * time.Second
	s, err := newSchedule(nominal, max, 2.0)
	if err != nil {
		t.Fatalf("newSchedule: %v", err)
	}

	expected := []struct {
		failed   bool
		returnD  time.Duration
		nextD    time.Duration
	}{
		{false, 5 * time.Second, 5 * time.Second},
		{true, 5 * time.Second, 10 * time.Second},
		{true, 10 * time.Second, 20 * time.Second},
		{true, 20 * time.Second, 40 * time.Second},
		{true, 40 * time.Second, 40 * time.Second}, // capped
		{false, 40 * time.Second, 5 * time.Second}, // recover
		{false, 5 * time.Second, 5 * time.Second},
	}

	for i, tc := range expected {
		got := s.next(tc.failed)
		if got != tc.returnD {
			t.Errorf("step %d: next() returned %v, want %v", i, got, tc.returnD)
		}
		if s.interval() != tc.nextD {
			t.Errorf("step %d: interval() = %v, want %v", i, s.interval(), tc.nextD)
		}
	}
}
