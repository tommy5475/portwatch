package daemon

import (
	"testing"
	"time"
)

func TestPauseInitiallyNotPaused(t *testing.T) {
	p := newPause()
	if p.IsPaused() {
		t.Fatal("expected not paused initially")
	}
}

func TestPauseTotalInitiallyZero(t *testing.T) {
	p := newPause()
	if p.TotalPaused() != 0 {
		t.Fatalf("expected zero total, got %v", p.TotalPaused())
	}
}

func TestPauseSetsPaused(t *testing.T) {
	p := newPause()
	p.Pause()
	if !p.IsPaused() {
		t.Fatal("expected paused after Pause()")
	}
}

func TestPauseResumeClears(t *testing.T) {
	p := newPause()
	p.Pause()
	p.Resume()
	if p.IsPaused() {
		t.Fatal("expected not paused after Resume()")
	}
}

func TestPauseIdempotent(t *testing.T) {
	p := newPause()
	p.Pause()
	p.Pause() // second call should be a no-op
	p.Resume()
	if p.IsPaused() {
		t.Fatal("expected not paused after Resume()")
	}
}

func TestResumeWithoutPauseIsNoop(t *testing.T) {
	p := newPause()
	p.Resume() // should not panic or change state
	if p.IsPaused() {
		t.Fatal("unexpected paused state")
	}
	if p.TotalPaused() != 0 {
		t.Fatalf("expected zero total after spurious Resume, got %v", p.TotalPaused())
	}
}

func TestPauseTotalAccumulates(t *testing.T) {
	p := newPause()
	p.Pause()
	time.Sleep(20 * time.Millisecond)
	p.Resume()
	if p.TotalPaused() < 10*time.Millisecond {
		t.Fatalf("expected accumulated pause >= 10ms, got %v", p.TotalPaused())
	}
}

func TestPauseTotalIncludesOngoingPause(t *testing.T) {
	p := newPause()
	p.Pause()
	time.Sleep(10 * time.Millisecond)
	if p.TotalPaused() < 5*time.Millisecond {
		t.Fatalf("expected TotalPaused to include ongoing pause, got %v", p.TotalPaused())
	}
}

func TestPauseReset(t *testing.T) {
	p := newPause()
	p.Pause()
	time.Sleep(10 * time.Millisecond)
	p.Reset()
	if p.IsPaused() {
		t.Fatal("expected not paused after Reset")
	}
	if p.TotalPaused() != 0 {
		t.Fatalf("expected zero total after Reset, got %v", p.TotalPaused())
	}
}
