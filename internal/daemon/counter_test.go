package daemon

import (
	"sync"
	"testing"
)

func TestCounterInitialValue(t *testing.T) {
	c := newCounter()
	if got := c.Get(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCounterInc(t *testing.T) {
	c := newCounter()
	if v := c.Inc(); v != 1 {
		t.Fatalf("expected 1 after first Inc, got %d", v)
	}
	if v := c.Inc(); v != 2 {
		t.Fatalf("expected 2 after second Inc, got %d", v)
	}
	if got := c.Get(); got != 2 {
		t.Fatalf("Get expected 2, got %d", got)
	}
}

func TestCounterAdd(t *testing.T) {
	c := newCounter()
	c.Add(10)
	if got := c.Get(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
	c.Add(-3)
	if got := c.Get(); got != 7 {
		t.Fatalf("expected 7 after negative Add, got %d", got)
	}
}

func TestCounterReset(t *testing.T) {
	c := newCounter()
	c.Add(42)
	prev := c.Reset()
	if prev != 42 {
		t.Fatalf("Reset should return previous value 42, got %d", prev)
	}
	if got := c.Get(); got != 0 {
		t.Fatalf("expected 0 after Reset, got %d", got)
	}
}

func TestCounterConcurrentInc(t *testing.T) {
	const goroutines = 100
	const incsEach = 50

	c := newCounter()
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incsEach; j++ {
				c.Inc()
			}
		}()
	}

	wg.Wait()

	expected := int64(goroutines * incsEach)
	if got := c.Get(); got != expected {
		t.Fatalf("expected %d after concurrent Inc, got %d", expected, got)
	}
}
