package daemon

import (
	"sync"
	"testing"
)

func TestTallyInitialValueIsZero(t *testing.T) {
	tl := newTally(100)
	if got := tl.Get("x"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestTallyIncReturnsNewValue(t *testing.T) {
	tl := newTally(100)
	if v := tl.Inc("a"); v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}
}

func TestTallyAddAccumulates(t *testing.T) {
	tl := newTally(1000)
	tl.Add("a", 10)
	tl.Add("a", 5)
	if got := tl.Get("a"); got != 15 {
		t.Fatalf("expected 15, got %d", got)
	}
}

func TestTallyCeiling(t *testing.T) {
	tl := newTally(5)
	tl.Add("a", 10)
	if got := tl.Get("a"); got != 5 {
		t.Fatalf("expected ceiling 5, got %d", got)
	}
}

func TestTallyNegativeAddClampedToZero(t *testing.T) {
	tl := newTally(100)
	tl.Add("a", -50)
	if got := tl.Get("a"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestTallyReset(t *testing.T) {
	tl := newTally(100)
	tl.Inc("a")
	tl.Reset("a")
	if got := tl.Get("a"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestTallySnapshot(t *testing.T) {
	tl := newTally(100)
	tl.Inc("a")
	tl.Add("b", 3)
	snap := tl.Snapshot()
	if snap["a"] != 1 || snap["b"] != 3 {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
}

func TestTallyLen(t *testing.T) {
	tl := newTally(100)
	tl.Inc("x")
	tl.Inc("y")
	if l := tl.Len(); l != 2 {
		t.Fatalf("expected len 2, got %d", l)
	}
}

func TestTallyDefaultCeilingOnInvalidArg(t *testing.T) {
	tl := newTally(0)
	for i := 0; i < 200; i++ {
		tl.Inc("k")
	}
	if tl.Get("k") != 200 {
		t.Fatalf("expected 200 with no ceiling, got %d", tl.Get("k"))
	}
}

func TestTallyConcurrentInc(t *testing.T) {
	tl := newTally(0)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tl.Inc("shared")
		}()
	}
	wg.Wait()
	if got := tl.Get("shared"); got != 100 {
		t.Fatalf("expected 100, got %d", got)
	}
}
