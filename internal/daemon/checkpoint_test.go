package daemon

import (
	"testing"
	"time"
)

func TestCheckpointInitialCountIsZero(t *testing.T) {
	cp := newCheckpoint()
	if got := cp.Count("boot"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCheckpointLastAtZeroBeforeMark(t *testing.T) {
	cp := newCheckpoint()
	if !cp.LastAt("boot").IsZero() {
		t.Fatal("expected zero time before first mark")
	}
}

func TestCheckpointFirstAtZeroBeforeMark(t *testing.T) {
	cp := newCheckpoint()
	if !cp.FirstAt("boot").IsZero() {
		t.Fatal("expected zero time before first mark")
	}
}

func TestCheckpointMarkIncrementsCount(t *testing.T) {
	cp := newCheckpoint()
	cp.Mark("scan")
	cp.Mark("scan")
	cp.Mark("scan")
	if got := cp.Count("scan"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestCheckpointLastAtUpdatesOnMark(t *testing.T) {
	cp := newCheckpoint()
	before := time.Now()
	cp.Mark("scan")
	after := time.Now()
	la := cp.LastAt("scan")
	if la.Before(before) || la.After(after) {
		t.Fatalf("LastAt %v not in [%v, %v]", la, before, after)
	}
}

func TestCheckpointFirstAtDoesNotChangeOnSubsequentMarks(t *testing.T) {
	cp := newCheckpoint()
	cp.Mark("scan")
	first := cp.FirstAt("scan")
	time.Sleep(2 * time.Millisecond)
	cp.Mark("scan")
	if got := cp.FirstAt("scan"); !got.Equal(first) {
		t.Fatalf("FirstAt changed: was %v, now %v", first, got)
	}
}

func TestCheckpointNamesReturnsMarked(t *testing.T) {
	cp := newCheckpoint()
	cp.Mark("a")
	cp.Mark("b")
	names := cp.Names()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}

func TestCheckpointResetClearsAll(t *testing.T) {
	cp := newCheckpoint()
	cp.Mark("scan")
	cp.Reset()
	if cp.Count("scan") != 0 {
		t.Fatal("expected count 0 after reset")
	}
	if len(cp.Names()) != 0 {
		t.Fatal("expected no names after reset")
	}
}

func TestCheckpointKeyIsolation(t *testing.T) {
	cp := newCheckpoint()
	cp.Mark("a")
	cp.Mark("a")
	cp.Mark("b")
	if cp.Count("a") != 2 {
		t.Fatalf("expected count 2 for 'a', got %d", cp.Count("a"))
	}
	if cp.Count("b") != 1 {
		t.Fatalf("expected count 1 for 'b', got %d", cp.Count("b"))
	}
}
