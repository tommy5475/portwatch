package daemon

import (
	"testing"
	"time"
)

func TestCooldownReadyWhenNoEntry(t *testing.T) {
	cd := newCooldown(10 * time.Second)
	if !cd.Ready("k") {
		t.Fatal("expected ready for unseen key")
	}
}

func TestCooldownNotReadyAfterMark(t *testing.T) {
	cd := newCooldown(10 * time.Second)
	cd.Mark("k")
	if cd.Ready("k") {
		t.Fatal("expected not ready immediately after mark")
	}
}

func TestCooldownReadyAfterPeriodElapses(t *testing.T) {
	now := time.Now()
	cd := newCooldown(5 * time.Second)
	cd.nowFn = func() time.Time { return now }
	cd.Mark("k")

	cd.nowFn = func() time.Time { return now.Add(6 * time.Second) }
	if !cd.Ready("k") {
		t.Fatal("expected ready after period elapsed")
	}
}

func TestCooldownRemainingIsZeroWhenReady(t *testing.T) {
	cd := newCooldown(5 * time.Second)
	if r := cd.Remaining("k"); r != 0 {
		t.Fatalf("expected 0 remaining for unseen key, got %v", r)
	}
}

func TestCooldownRemainingDecreasesOverTime(t *testing.T) {
	now := time.Now()
	cd := newCooldown(10 * time.Second)
	cd.nowFn = func() time.Time { return now }
	cd.Mark("k")

	cd.nowFn = func() time.Time { return now.Add(4 * time.Second) }
	r := cd.Remaining("k")
	if r != 6*time.Second {
		t.Fatalf("expected 6s remaining, got %v", r)
	}
}

func TestCooldownResetAllowsImmediateReady(t *testing.T) {
	cd := newCooldown(10 * time.Second)
	cd.Mark("k")
	cd.Reset("k")
	if !cd.Ready("k") {
		t.Fatal("expected ready after reset")
	}
}

func TestCooldownKeyIsolation(t *testing.T) {
	cd := newCooldown(10 * time.Second)
	cd.Mark("a")
	if !cd.Ready("b") {
		t.Fatal("marking 'a' should not affect 'b'")
	}
}

func TestCooldownDefaultsOnInvalidPeriod(t *testing.T) {
	cd := newCooldown(-1 * time.Second)
	if cd.period != 5*time.Second {
		t.Fatalf("expected default period 5s, got %v", cd.period)
	}
}
