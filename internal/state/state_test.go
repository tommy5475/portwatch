package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"portwatch/internal/state"
)

func tmpStore(t *testing.T) *state.Store {
	t.Helper()
	dir := t.TempDir()
	return state.New(filepath.Join(dir, "state.json"))
}

func TestLoadEmpty(t *testing.T) {
	store := tmpStore(t)
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %d", len(snap.Ports))
	}
}

func TestSaveAndLoad(t *testing.T) {
	store := tmpStore(t)
	original := state.Snapshot{
		Ports: []state.PortEntry{
			{Port: 80, Protocol: "tcp", Address: "0.0.0.0"},
			{Port: 443, Protocol: "tcp", Address: "0.0.0.0"},
		},
	}
	if err := store.Save(original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Ports) != len(original.Ports) {
		t.Errorf("expected %d ports, got %d", len(original.Ports), len(loaded.Ports))
	}
	if loaded.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set after Save")
	}
}

func TestSaveCreatesFileWithRestrictedPermissions(t *testing.T) {
	store := tmpStore(t)
	if err := store.Save(state.Snapshot{}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	s2 := state.New(path)
	_ = s2.Save(state.Snapshot{})
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("expected file mode 0600, got %v", info.Mode().Perm())
	}
}

func TestDiff(t *testing.T) {
	prev := state.Snapshot{
		Ports: []state.PortEntry{
			{Port: 22, Protocol: "tcp", Address: "0.0.0.0"},
			{Port: 80, Protocol: "tcp", Address: "0.0.0.0"},
		},
	}
	curr := state.Snapshot{
		Ports: []state.PortEntry{
			{Port: 80, Protocol: "tcp", Address: "0.0.0.0"},
			{Port: 8080, Protocol: "tcp", Address: "0.0.0.0"},
		},
	}
	opened, closed := state.Diff(prev, curr)
	if len(opened) != 1 || opened[0].Port != 8080 {
		t.Errorf("expected port 8080 opened, got %v", opened)
	}
	if len(closed) != 1 || closed[0].Port != 22 {
		t.Errorf("expected port 22 closed, got %v", closed)
	}
}

func TestDiffNoChanges(t *testing.T) {
	ports := []state.PortEntry{{Port: 443, Protocol: "tcp", Address: "0.0.0.0"}}
	snap := state.Snapshot{Ports: ports}
	opened, closed := state.Diff(snap, snap)
	if len(opened) != 0 || len(closed) != 0 {
		t.Errorf("expected no changes, got opened=%v closed=%v", opened, closed)
	}
}
