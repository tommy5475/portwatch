package daemon

import (
	"testing"
	"time"
)

func TestHeartbeatInitialLastIsZero(t *testing.T) {
	hb := newHeartbeat()
	if !hb.last().IsZero() {
		t.Fatal("expected zero time before first pulse")
	}
}

func TestHeartbeatInitialCountIsZero(t *testing.T) {
	hb := newHeartbeat()
	if hb.count() != 0 {
		t.Fatalf("expected count 0, got %d", hb.count())
	}
}

func TestHeartbeatPulseIncrementsCount(t *testing.T) {
	hb := newHeartbeat()
	hb.pulse()
	hb.pulse()
	if hb.count() != 2 {
		t.Fatalf("expected count 2, got %d", hb.count())
	}
}

func TestHeartbeatLastIsSetAfterPulse(t *testing.T) {
	hb := newHeartbeat()
	before := time.Now()
	hb.pulse()
	after := time.Now()
	l := hb.last()
	if l.Before(before) || l.After(after) {
		t.Fatalf("last pulse time %v not in expected range [%v, %v]", l, before, after)
	}
}

func TestHeartbeatAgeIsPositiveAfterPulse(t *testing.T) {
	hb := newHeartbeat()
	hb.pulse()
	if hb.age() < 0 {
		t.Fatal("age should be non-negative after pulse")
	}
}

func TestHeartbeatAgeWithoutPulseEqualsUptime(t *testing.T) {
	hb := newHeartbeat()
	time.Sleep(2 * time.Millisecond)
	age := hb.age()
	if age < time.Millisecond {
		t.Fatalf("expected age >= 1ms before first pulse, got %v", age)
	}
}

func TestHeartbeatAliveWithinThreshold(t *testing.T) {
	hb := newHeartbeat()
	hb.pulse()
	if !hb.alive(time.Second) {
		t.Fatal("expected alive=true immediately after pulse")
	}
}

func TestHeartbeatNotAliveAfterThresholdExceeded(t *testing.T) {
	hb := newHeartbeat()
	hb.pulse()
	time.Sleep(10 * time.Millisecond)
	if hb.alive(time.Nanosecond) {
		t.Fatal("expected alive=false after threshold exceeded")
	}
}

func TestHeartbeatConcurrentPulse(t *testing.T) {
	hb := newHeartbeat()
	const goroutines = 50
	done := make(chan struct{}, goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			hb.pulse()
			done <- struct{}{}
		}()
	}
	for i := 0; i < goroutines; i++ {
		<-done
	}
	if hb.count() != goroutines {
		t.Fatalf("expected count %d, got %d", goroutines, hb.count())
	}
}
