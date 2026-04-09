package daemon

import (
	"testing"
	"time"
)

func TestThrottleAllowsFirstCall(t *testing.T) {
	th := newThrottle(5 * time.Second)
	if !th.Allow("tcp:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestThrottleBlocksWithinCooldown(t *testing.T) {
	th := newThrottle(5 * time.Second)
	th.Allow("tcp:8080")
	if th.Allow("tcp:8080") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestThrottleAllowsAfterCooldown(t *testing.T) {
	now := time.Unix(1_000_000, 0)
	th := newThrottle(5 * time.Second)
	th.now = func() time.Time { return now }

	th.Allow("tcp:8080")

	th.now = func() time.Time { return now.Add(6 * time.Second) }
	if !th.Allow("tcp:8080") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestThrottleIndependentKeys(t *testing.T) {
	th := newThrottle(5 * time.Second)
	th.Allow("tcp:8080")
	if !th.Allow("tcp:9090") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestThrottleReset(t *testing.T) {
	th := newThrottle(5 * time.Second)
	th.Allow("tcp:8080")
	th.Reset("tcp:8080")
	if !th.Allow("tcp:8080") {
		t.Fatal("expected allow after reset")
	}
}

func TestThrottleDefaultCooldownOnInvalidArg(t *testing.T) {
	th := newThrottle(0)
	if th.cooldown != 30*time.Second {
		t.Fatalf("expected default 30s cooldown, got %v", th.cooldown)
	}
}

func TestThrottleLen(t *testing.T) {
	th := newThrottle(5 * time.Second)
	th.Allow("a")
	th.Allow("b")
	th.Allow("c")
	if th.Len() != 3 {
		t.Fatalf("expected Len 3, got %d", th.Len())
	}
	th.Reset("b")
	if th.Len() != 2 {
		t.Fatalf("expected Len 2 after reset, got %d", th.Len())
	}
}
