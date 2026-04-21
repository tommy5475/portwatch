package daemon

import (
	"sync"
	"time"
)

// sampler collects periodic numeric observations and exposes
// summary statistics: min, max, mean, and the last recorded value.
type sampler struct {
	mu      sync.Mutex
	min     float64
	max     float64
	sum     float64
	count   int64
	last    float64
	lastAt  time.Time
	startAt time.Time
}

func newSampler() *sampler {
	return &sampler{startAt: time.Now()}
}

// record adds a new observation.
func (s *sampler) record(v float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.count == 0 || v < s.min {
		s.min = v
	}
	if s.count == 0 || v > s.max {
		s.max = v
	}
	s.sum += v
	s.count++
	s.last = v
	s.lastAt = time.Now()
}

// mean returns the arithmetic mean, or 0 if no observations.
func (s *sampler) mean() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.count == 0 {
		return 0
	}
	return s.sum / float64(s.count)
}

// snapshot returns a point-in-time copy of all statistics.
func (s *sampler) stats() samplerStats {
	s.mu.Lock()
	defer s.mu.Unlock()
	return samplerStats{
		Min:    s.min,
		Max:    s.max,
		Mean:   func() float64 {
			if s.count == 0 {
				return 0
			}
			return s.sum / float64(s.count)
		}(),
		Last:   s.last,
		LastAt: s.lastAt,
		Count:  s.count,
		Uptime: time.Since(s.startAt),
	}
}

// reset clears all observations.
func (s *sampler) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.min = 0
	s.max = 0
	s.sum = 0
	s.count = 0
	s.last = 0
	s.lastAt = time.Time{}
	s.startAt = time.Now()
}

type samplerStats struct {
	Min    float64
	Max    float64
	Mean   float64
	Last   float64
	LastAt time.Time
	Count  int64
	Uptime time.Duration
}
