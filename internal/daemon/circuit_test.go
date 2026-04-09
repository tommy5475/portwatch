package daemon

import (
	"testing"
	"time"
)

func TestCircuitInitiallyClosed(t *testing.T) {
	c := newCircuit(3, time.Second)
	if !c.allow() {
		t.Fatal("expected circuit to allow calls when closed")
	}
}

func TestCircuitOpensAfterThreshold(t *testing.T) {
	c := newCircuit(3, time.Minute)
	for i := 0; i < 3; i++ {
		c.recordFailure()
	}
	if c.allow() {
		t.Fatal("expected circuit to block calls when open")
	}
}

func TestCircuitDoesNotOpenBeforeThreshold(t *testing.T) {
	c := newCircuit(3, time.Minute)
	c.recordFailure()
	c.recordFailure()
	if !c.allow() {
		t.Fatal("expected circuit to remain closed below threshold")
	}
}

func TestCircuitHalfOpenAfterCooldown(t *testing.T) {
	c := newCircuit(1, time.Millisecond)
	c.recordFailure()
	time.Sleep(5 * time.Millisecond)
	if !c.allow() {
		t.Fatal("expected circuit to allow probe after cooldown")
	}
	state, _, _ := c.stats()
	if state != circuitHalfOpen {
		t.Fatalf("expected half-open state, got %d", state)
	}
}

func TestCircuitClosesOnSuccessFromHalfOpen(t *testing.T) {
	c := newCircuit(1, time.Millisecond)
	c.recordFailure()
	time.Sleep(5 * time.Millisecond)
	c.allow() // transitions to half-open
	c.recordSuccess()
	state, failures, _ := c.stats()
	if state != circuitClosed {
		t.Fatalf("expected closed state after success, got %d", state)
	}
	if failures != 0 {
		t.Fatalf("expected zero failures after success, got %d", failures)
	}
}

func TestCircuitTripCountIncrements(t *testing.T) {
	c := newCircuit(2, time.Millisecond)
	for round := 0; round < 3; round++ {
		c.recordFailure()
		c.recordFailure()
		time.Sleep(5 * time.Millisecond)
		c.allow()
		c.recordSuccess()
	}
	_, _, trips := c.stats()
	if trips != 3 {
		t.Fatalf("expected 3 trips, got %d", trips)
	}
}

func TestCircuitDefaultsOnInvalidArgs(t *testing.T) {
	c := newCircuit(0, 0)
	if c.threshold != 3 {
		t.Fatalf("expected default threshold 3, got %d", c.threshold)
	}
	if c.cooldown != 30*time.Second {
		t.Fatalf("expected default cooldown 30s, got %v", c.cooldown)
	}
}
