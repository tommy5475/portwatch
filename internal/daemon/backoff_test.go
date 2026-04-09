package daemon

import (
	"testing"
	"time"
)

func TestBackoffInitialState(t *testing.T) {
	b := newBackoff(2*time.Second, 60*time.Second)
	if b.Failures() != 0 {
		t.Fatalf("expected 0 failures, got %d", b.Failures())
	}
	if b.Current() != 2*time.Second {
		t.Fatalf("expected initial current=2s, got %v", b.Current())
	}
}

func TestBackoffFailureDoublesDelay(t *testing.T) {
	b := newBackoff(2*time.Second, 60*time.Second)

	d1 := b.Failure() // 2s
	d2 := b.Failure() // 4s
	d3 := b.Failure() // 8s

	if d1 != 2*time.Second {
		t.Errorf("first failure: want 2s, got %v", d1)
	}
	if d2 != 4*time.Second {
		t.Errorf("second failure: want 4s, got %v", d2)
	}
	if d3 != 8*time.Second {
		t.Errorf("third failure: want 8s, got %v", d3)
	}
}

func TestBackoffCapsAtMaxDelay(t *testing.T) {
	b := newBackoff(2*time.Second, 8*time.Second)

	for i := 0; i < 10; i++ {
		d := b.Failure()
		if d > 8*time.Second {
			t.Fatalf("delay %v exceeded maxDelay 8s on iteration %d", d, i)
		}
	}
}

func TestBackoffSuccessResetsState(t *testing.T) {
	b := newBackoff(2*time.Second, 60*time.Second)

	b.Failure()
	b.Failure()
	b.Failure()

	if b.Failures() != 3 {
		t.Fatalf("expected 3 failures before reset, got %d", b.Failures())
	}

	b.Success()

	if b.Failures() != 0 {
		t.Errorf("expected 0 failures after Success(), got %d", b.Failures())
	}
	if b.Current() != 2*time.Second {
		t.Errorf("expected current reset to 2s, got %v", b.Current())
	}

	// Next failure should start from base again.
	d := b.Failure()
	if d != 2*time.Second {
		t.Errorf("after reset, first failure should be 2s, got %v", d)
	}
}

func TestBackoffDefaultsOnInvalidArgs(t *testing.T) {
	b := newBackoff(0, 0)
	if b.BaseDelay <= 0 {
		t.Error("BaseDelay should be positive when zero is passed")
	}
	if b.MaxDelay < b.BaseDelay {
		t.Error("MaxDelay should be >= BaseDelay")
	}
}
