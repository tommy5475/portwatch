// Package state manages persistence and comparison of observed port sets.
package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// PortEntry represents a single observed open port.
type PortEntry struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

// Snapshot holds a complete set of observed open ports at a point in time.
type Snapshot struct {
	Ports []PortEntry `json:"ports"`
}

// Diff describes ports that have opened or closed between two snapshots.
type Diff struct {
	Opened []PortEntry
	Closed []PortEntry
}

// IsEmpty reports whether the diff contains no changes.
func (d Diff) IsEmpty() bool {
	return len(d.Opened) == 0 && len(d.Closed) == 0
}

// Store persists snapshots to disk.
type Store struct {
	path string
}

// New creates a Store that persists state to the given file path.
func New(path string) *Store {
	return &Store{path: path}
}

// Load reads the last saved snapshot from disk.
// If no file exists an empty Snapshot is returned without error.
func (s *Store) Load() (Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Save writes the snapshot to disk with restricted permissions.
func (s *Store) Save(snap Snapshot) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

// Diff computes the difference between two snapshots.
func Diff(prev, curr Snapshot) Diff {
	prevSet := toSet(prev)
	currSet := toSet(curr)
	var d Diff
	for k, e := range currSet {
		if _, ok := prevSet[k]; !ok {
			d.Opened = append(d.Opened, e)
		}
	}
	for k, e := range prevSet {
		if _, ok := currSet[k]; !ok {
			d.Closed = append(d.Closed, e)
		}
	}
	return d
}

func toSet(snap Snapshot) map[string]PortEntry {
	m := make(map[string]PortEntry, len(snap.Ports))
	for _, p := range snap.Ports {
		key := p.Protocol + ":" + string(rune(p.Port))
		m[key] = p
	}
	return m
}
