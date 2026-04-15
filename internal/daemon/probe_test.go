package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestProbeInitiallyLiveNotReady(t *testing.T) {
	p := newProbe()
	if !p.IsLive() {
		t.Fatal("expected probe to be live initially")
	}
	if p.IsReady() {
		t.Fatal("expected probe to be not ready initially")
	}
}

func TestProbeMarkReady(t *testing.T) {
	p := newProbe()
	p.MarkReady()
	if !p.IsReady() {
		t.Fatal("expected probe to be ready after MarkReady")
	}
}

func TestProbeMarkNotReady(t *testing.T) {
	p := newProbe()
	p.MarkReady()
	p.MarkNotReady()
	if p.IsReady() {
		t.Fatal("expected probe to be not ready after MarkNotReady")
	}
}

func TestProbeMarkNotLiveIncrementsFailCount(t *testing.T) {
	p := newProbe()
	p.MarkNotLive()
	p.MarkNotLive()
	if got := p.FailCount(); got != 2 {
		t.Fatalf("expected fail count 2, got %d", got)
	}
}

func TestProbeMarkNotLiveSetsFailedAt(t *testing.T) {
	p := newProbe()
	before := time.Now()
	p.MarkNotLive()
	after := time.Now()
	ft := p.FailedAt()
	if ft.Before(before) || ft.After(after) {
		t.Fatalf("failedAt %v not between %v and %v", ft, before, after)
	}
}

func TestProbeMarkLiveRestoresLiveness(t *testing.T) {
	p := newProbe()
	p.MarkNotLive()
	if p.IsLive() {
		t.Fatal("expected probe not live after MarkNotLive")
	}
	p.MarkLive()
	if !p.IsLive() {
		t.Fatal("expected probe live after MarkLive")
	}
}

func TestProbeResetClearsState(t *testing.T) {
	p := newProbe()
	p.MarkNotLive()
	p.MarkReady()
	p.Reset()
	if !p.IsLive() {
		t.Fatal("expected live after reset")
	}
	if p.IsReady() {
		t.Fatal("expected not ready after reset")
	}
	if p.FailCount() != 0 {
		t.Fatalf("expected fail count 0 after reset, got %d", p.FailCount())
	}
	if !p.FailedAt().IsZero() {
		t.Fatal("expected zero failedAt after reset")
	}
}

func TestProbeConcurrentAccess(t *testing.T) {
	p := newProbe()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				p.MarkNotLive()
			} else {
				p.MarkLive()
			}
			_ = p.IsLive()
			_ = p.FailCount()
		}(i)
	}
	wg.Wait()
}
