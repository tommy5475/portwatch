package daemon

import (
	"sync"
	"testing"
)

func TestGaugeInitialValueIsZero(t *testing.T) {
	g := newGauge()
	if got := g.Get(); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestGaugeCountInitiallyZero(t *testing.T) {
	g := newGauge()
	if g.Count() != 0 {
		t.Fatalf("expected count 0, got %d", g.Count())
	}
}

func TestGaugeSetUpdatesValue(t *testing.T) {
	g := newGauge()
	g.Set(3.14)
	if got := g.Get(); got != 3.14 {
		t.Fatalf("expected 3.14, got %v", got)
	}
}

func TestGaugeMinMaxSingleValue(t *testing.T) {
	g := newGauge()
	g.Set(7.5)
	if g.Min() != 7.5 {
		t.Fatalf("expected min 7.5, got %v", g.Min())
	}
	if g.Max() != 7.5 {
		t.Fatalf("expected max 7.5, got %v", g.Max())
	}
}

func TestGaugeTracksMinMax(t *testing.T) {
	g := newGauge()
	g.Set(5.0)
	g.Set(1.0)
	g.Set(9.0)
	if g.Min() != 1.0 {
		t.Fatalf("expected min 1.0, got %v", g.Min())
	}
	if g.Max() != 9.0 {
		t.Fatalf("expected max 9.0, got %v", g.Max())
	}
}

func TestGaugeMinIsZeroBeforeSet(t *testing.T) {
	g := newGauge()
	if g.Min() != 0 {
		t.Fatalf("expected min 0 before any Set, got %v", g.Min())
	}
}

func TestGaugeCountIncrementsOnSet(t *testing.T) {
	g := newGauge()
	g.Set(1)
	g.Set(2)
	g.Set(3)
	if g.Count() != 3 {
		t.Fatalf("expected count 3, got %d", g.Count())
	}
}

func TestGaugeResetClearsState(t *testing.T) {
	g := newGauge()
	g.Set(42.0)
	g.Reset()
	if g.Get() != 0 {
		t.Fatalf("expected 0 after reset, got %v", g.Get())
	}
	if g.Count() != 0 {
		t.Fatalf("expected count 0 after reset, got %d", g.Count())
	}
	if g.Min() != 0 {
		t.Fatalf("expected min 0 after reset, got %v", g.Min())
	}
}

func TestGaugeConcurrentSet(t *testing.T) {
	g := newGauge()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		v := float64(i)
		go func() {
			defer wg.Done()
			g.Set(v)
		}()
	}
	wg.Wait()
	if g.Count() != 100 {
		t.Fatalf("expected count 100, got %d", g.Count())
	}
	if g.Max() < 99.0 {
		t.Fatalf("expected max >= 99, got %v", g.Max())
	}
}
