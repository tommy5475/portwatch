package daemon

import (
	"testing"
	"time"
)

func TestBudgetInitiallyFull(t *testing.T) {
	b := newBudget(5, time.Second)
	if got := b.Remaining(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestBudgetSpendReducesTokens(t *testing.T) {
	b := newBudget(5, time.Second)
	if !b.Spend(3) {
		t.Fatal("expected spend to succeed")
	}
	if got := b.Remaining(); got != 2 {
		t.Fatalf("expected 2 remaining, got %d", got)
	}
}

func TestBudgetSpendFailsWhenExhausted(t *testing.T) {
	b := newBudget(3, time.Second)
	b.Spend(3)
	if b.Spend(1) {
		t.Fatal("expected spend to fail when budget exhausted")
	}
}

func TestBudgetResetRestoresFull(t *testing.T) {
	b := newBudget(4, time.Second)
	b.Spend(4)
	b.Reset()
	if got := b.Remaining(); got != 4 {
		t.Fatalf("expected 4 after reset, got %d", got)
	}
}

func TestBudgetSpendZeroAlwaysSucceeds(t *testing.T) {
	b := newBudget(1, time.Second)
	b.Spend(1) // exhaust
	if !b.Spend(0) {
		t.Fatal("spending 0 tokens should always succeed")
	}
}

func TestBudgetReplenishesOverTime(t *testing.T) {
	b := newBudget(10, 10*time.Millisecond)
	b.Spend(10) // exhaust
	time.Sleep(35 * time.Millisecond)
	remaining := b.Remaining()
	if remaining < 3 {
		t.Fatalf("expected at least 3 tokens replenished, got %d", remaining)
	}
}

func TestBudgetDoesNotExceedCap(t *testing.T) {
	b := newBudget(5, 10*time.Millisecond)
	// already full; sleeping should not exceed cap
	time.Sleep(60 * time.Millisecond)
	if got := b.Remaining(); got > 5 {
		t.Fatalf("tokens exceeded cap: got %d", got)
	}
}

func TestBudgetDefaultsOnInvalidArgs(t *testing.T) {
	b := newBudget(0, 0)
	if b.cap <= 0 {
		t.Fatal("expected positive default cap")
	}
	if b.rate <= 0 {
		t.Fatal("expected positive default rate")
	}
}
