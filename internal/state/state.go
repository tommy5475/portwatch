// Package state provides persistent storage for port scan snapshots,
// enabling portwatch to detect changes between daemon runs.
package state

import (
	"encoding/json"
	"os"
	"time"
)

// PortEntry represents a single observed open port.
type PortEntry struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
}

// Snapshot holds the full port state captured at a point in time.
type Snapshot struct {
	CapturedAt time.Time   `json:"captured_at"`
	Ports      []PortEntry `json:"ports"`
}

// Store manages reading and writing state snapshots to disk.
type Store struct {
	path string
}

// New creates a new Store that persists state to the given file path.
func New(path string) *Store {
	return &Store{path: path}
}

// Save writes the snapshot to disk, overwriting any previous state.
func (s *Store) Save(snap Snapshot) error {
	snap.CapturedAt = time.Now()
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

// Load reads the last persisted snapshot from disk.
// Returns an empty Snapshot and no error when the file does not yet exist.
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

// Diff compares two snapshots and returns the ports that were opened or closed.
func Diff(prev, curr Snapshot) (opened, closed []PortEntry) {
	prevSet := toSet(prev.Ports)
	currSet := toSet(curr.Ports)

	for key, entry := range currSet {
		if _, exists := prevSet[key]; !exists {
			opened = append(opened, entry)
		}
	}
	for key, entry := range prevSet {
		if _, exists := currSet[key]; !exists {
			closed = append(closed, entry)
		}
	}
	return opened, closed
}

func toSet(ports []PortEntry) map[string]PortEntry {
	m := make(map[string]PortEntry, len(ports))
	for _, p := range ports {
		key := p.Protocol + ":" + p.Address + ":" + string(rune(p.Port+'0'))
		m[key] = p
	}
	return m
}
