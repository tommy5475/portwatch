package daemon

import (
	"testing"
	"time"
)

func TestEpochInitialGenerationIsZero(t *testing.T) {
	e := newEpoch()
	if e.generation() != 0 {
		t.Fatalf("expected initial generation 0, got %d", e.generation())
	}
}

func TestEpochAdvanceIncrementsGeneration(t *testing.T) {
	e := newEpoch()
	e.advance()
	if e.generation() != 1 {
		t.Fatalf("expected generation 1, got %d", e.generation())
	}
	e.advance()
	if e.generation() != 2 {
		t.Fatalf("expected generation 2, got %d", e.generation())
	}
}

func TestEpochLastAdvanceUpdatesOnAdvance(t *testing.T) {
	e := newEpoch()
	before := time.Now()
	e.advance()
	after := time.Now()
	l := e.lastAdvance()
	if l.Before(before) || l.After(after) {
		t.Fatalf("lastAdvance %v not in [%v, %v]", l, before, after)
	}
}

func TestEpochAgeIsPositive(t *testing.T) {
	e := newEpoch()
	time.Sleep(time.Millisecond)
	if e.age() <= 0 {
		t.Fatal("expected positive age")
	}
}

func TestEpochGenerationIsMonotonic(t *testing.T) {
	e := newEpoch()
	const n = 50
	for i := 0; i < n; i++ {
		e.advance()
	}
	if e.generation() != n {
		t.Fatalf("expected generation %d, got %d", n, e.generation())
	}
}

func TestEpochConcurrentAdvance(t *testing.T) {
	e := newEpoch()
	const workers = 10
	const advances = 100
	done := make(chan struct{})
	for i := 0; i < workers; i++ {
		go func() {
			for j := 0; j < advances; j++ {
				e.advance()
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < workers; i++ {
		<-done
	}
	expected := uint64(workers * advances)
	if e.generation() != expected {
		t.Fatalf("expected generation %d, got %d", expected, e.generation())
	}
}
