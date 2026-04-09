package daemon

import (
	"testing"
	"time"

	"portwatch/internal/state"
)

func TestSnapshotInitialAge(t *testing.T) {
	s := newSnapshot()
	if s.age() != -1 {
		t.Errorf("expected age -1 before any update, got %v", s.age())
	}
}

func TestSnapshotInitialCount(t *testing.T) {
	s := newSnapshot()
	if s.count() != 0 {
		t.Errorf("expected count 0, got %d", s.count())
	}
}

func TestSnapshotUpdateIncrementsCount(t *testing.T) {
	s := newSnapshot()
	pm := state.PortMap{}
	s.update(pm)
	s.update(pm)
	if s.count() != 2 {
		t.Errorf("expected count 2, got %d", s.count())
	}
}

func TestSnapshotGetReturnsCopy(t *testing.T) {
	s := newSnapshot()
	pm := state.PortMap{"tcp:80": {Port: 80, Protocol: "tcp", Open: true}}
	s.update(pm)

	got, ts := s.get()
	if ts.IsZero() {
		t.Error("expected non-zero capture time")
	}
	if len(got) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got))
	}

	// Mutating the returned map must not affect the stored snapshot.
	delete(got, "tcp:80")
	got2, _ := s.get()
	if len(got2) != 1 {
		t.Error("snapshot was mutated through returned copy")
	}
}

func TestSnapshotAgeIsPositiveAfterUpdate(t *testing.T) {
	s := newSnapshot()
	s.update(state.PortMap{})
	time.Sleep(2 * time.Millisecond)
	if s.age() <= 0 {
		t.Errorf("expected positive age, got %v", s.age())
	}
}

func TestSnapshotUpdateReplacesData(t *testing.T) {
	s := newSnapshot()
	s.update(state.PortMap{"tcp:80": {Port: 80, Protocol: "tcp", Open: true}})
	s.update(state.PortMap{"tcp:443": {Port: 443, Protocol: "tcp", Open: true}})

	got, _ := s.get()
	if _, ok := got["tcp:80"]; ok {
		t.Error("old entry should have been replaced")
	}
	if _, ok := got["tcp:443"]; !ok {
		t.Error("new entry should be present")
	}
}
