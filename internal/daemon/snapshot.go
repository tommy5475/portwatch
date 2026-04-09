package daemon

import (
	"sync"
	"time"

	"portwatch/internal/state"
)

// snapshot holds the most recently observed port state along with
// metadata about when it was captured.
type snapshot struct {
	mu        sync.RWMutex
	ports     state.PortMap
	capturedAt time.Time
	scanCount  int64
}

func newSnapshot() *snapshot {
	return &snapshot{}
}

// update atomically replaces the current snapshot with a new port map.
func (s *snapshot) update(ports state.PortMap) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ports = ports
	s.capturedAt = time.Now()
	s.scanCount++
}

// get returns a shallow copy of the current port map and capture time.
func (s *snapshot) get() (state.PortMap, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(state.PortMap, len(s.ports))
	for k, v := range s.ports {
		copy[k] = v
	}
	return copy, s.capturedAt
}

// count returns the total number of scans recorded.
func (s *snapshot) count() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.scanCount
}

// age returns the duration since the snapshot was last updated.
// Returns -1 if no snapshot has been captured yet.
func (s *snapshot) age() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.capturedAt.IsZero() {
		return -1
	}
	return time.Since(s.capturedAt)
}
