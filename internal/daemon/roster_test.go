package daemon

import (
	"testing"
	"time"
)

func TestRosterJoinRegisters(t *testing.T) {
	r := newRoster()
	r.join("w1")
	if r.len() != 1 {
		t.Fatalf("expected 1 entry, got %d", r.len())
	}
}

func TestRosterIsAliveAfterJoin(t *testing.T) {
	r := newRoster()
	r.join("w1")
	if !r.isAlive("w1") {
		t.Fatal("expected worker to be alive after join")
	}
}

func TestRosterUnknownWorkerNotAlive(t *testing.T) {
	r := newRoster()
	if r.isAlive("ghost") {
		t.Fatal("unknown worker should not be alive")
	}
}

func TestRosterCheckinReturnsTrueForKnown(t *testing.T) {
	r := newRoster()
	r.join("w1")
	if !r.checkin("w1") {
		t.Fatal("checkin should return true for registered worker")
	}
}

func TestRosterCheckinReturnsFalseForUnknown(t *testing.T) {
	r := newRoster()
	if r.checkin("ghost") {
		t.Fatal("checkin should return false for unregistered worker")
	}
}

func TestRosterLeaveRemovesWorker(t *testing.T) {
	r := newRoster()
	r.join("w1")
	r.leave("w1")
	if r.len() != 0 {
		t.Fatalf("expected 0 entries after leave, got %d", r.len())
	}
}

func TestRosterSnapshotReturnsCopy(t *testing.T) {
	r := newRoster()
	r.join("w1")
	r.join("w2")
	snap := r.snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 snapshot entries, got %d", len(snap))
	}
}

func TestRosterMarkStaleTransitionsWorker(t *testing.T) {
	r := newRoster()
	r.join("w1")
	// back-date lastSeen so the worker appears stale
	r.mu.Lock()
	r.entries["w1"].lastSeen = time.Now().Add(-2 * time.Minute)
	r.mu.Unlock()

	stale := r.markStale(30 * time.Second)
	if stale != 1 {
		t.Fatalf("expected 1 stale worker, got %d", stale)
	}
	if r.isAlive("w1") {
		t.Fatal("worker should not be alive after markStale")
	}
}

func TestRosterMarkStaleSkipsFreshWorker(t *testing.T) {
	r := newRoster()
	r.join("w1")
	stale := r.markStale(30 * time.Second)
	if stale != 0 {
		t.Fatalf("expected 0 stale workers, got %d", stale)
	}
	if !r.isAlive("w1") {
		t.Fatal("fresh worker should remain alive")
	}
}

func TestRosterCheckinRestoresAlive(t *testing.T) {
	r := newRoster()
	r.join("w1")
	r.mu.Lock()
	r.entries["w1"].lastSeen = time.Now().Add(-2 * time.Minute)
	r.mu.Unlock()
	r.markStale(30 * time.Second)

	r.checkin("w1")
	if !r.isAlive("w1") {
		t.Fatal("worker should be alive after checkin")
	}
}
