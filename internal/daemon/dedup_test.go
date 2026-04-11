package daemon

import (
	"testing"
	"time"
)

func TestDedupFirstCallNotDuplicate(t *testing.T) {
	d := newDedup(time.Second)
	if d.IsDuplicate("port:8080") {
		t.Fatal("expected first call to not be a duplicate")
	}
}

func TestDedupSecondCallWithinWindowIsDuplicate(t *testing.T) {
	d := newDedup(time.Second)
	d.IsDuplicate("port:8080")
	if !d.IsDuplicate("port:8080") {
		t.Fatal("expected second call within window to be a duplicate")
	}
}

func TestDedupAfterWindowExpires(t *testing.T) {
	now := time.Unix(1000, 0)
	d := newDedup(time.Second)
	d.nowFn = func() time.Time { return now }

	d.IsDuplicate("port:9090")

	// advance past the window
	d.nowFn = func() time.Time { return now.Add(2 * time.Second) }
	if d.IsDuplicate("port:9090") {
		t.Fatal("expected call after window expiry to not be a duplicate")
	}
}

func TestDedupKeyIsolation(t *testing.T) {
	d := newDedup(time.Second)
	d.IsDuplicate("port:80")
	if d.IsDuplicate("port:443") {
		t.Fatal("expected different key to not be a duplicate")
	}
}

func TestDedupEvictRemovesExpiredKeys(t *testing.T) {
	now := time.Unix(2000, 0)
	d := newDedup(time.Second)
	d.nowFn = func() time.Time { return now }

	d.IsDuplicate("a")
	d.IsDuplicate("b")

	d.nowFn = func() time.Time { return now.Add(2 * time.Second) }
	d.Evict()

	if d.Len() != 0 {
		t.Fatalf("expected 0 keys after eviction, got %d", d.Len())
	}
}

func TestDedupReset(t *testing.T) {
	d := newDedup(time.Second)
	d.IsDuplicate("x")
	d.IsDuplicate("y")
	d.Reset()
	if d.Len() != 0 {
		t.Fatalf("expected 0 keys after reset, got %d", d.Len())
	}
}

func TestDedupDefaultsOnInvalidWindow(t *testing.T) {
	d := newDedup(-1)
	if d.window <= 0 {
		t.Fatal("expected positive default window")
	}
}
