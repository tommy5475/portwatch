package daemon

import (
	"errors"
	"time"
)

// schedule defines a variable-interval tick strategy that can back off
// during degraded states and recover to the nominal interval on success.
type schedule struct {
	nominal  time.Duration
	current  time.Duration
	maxDelay time.Duration
	factor   float64
}

func newSchedule(nominal, maxDelay time.Duration, factor float64) (*schedule, error) {
	if nominal <= 0 {
		return nil, errors.New("schedule: nominal interval must be positive")
	}
	if maxDelay < nominal {
		maxDelay = nominal
	}
	if factor < 1.0 {
		factor = 1.5
	}
	return &schedule{
		nominal:  nominal,
		current:  nominal,
		maxDelay: maxDelay,
		factor:   factor,
	}, nil
}

// next returns the current interval and advances the schedule on failure.
func (s *schedule) next(failed bool) time.Duration {
	d := s.current
	if failed {
		next := time.Duration(float64(s.current) * s.factor)
		if next > s.maxDelay {
			next = s.maxDelay
		}
		s.current = next
	} else {
		s.current = s.nominal
	}
	return d
}

// reset restores the schedule to its nominal interval.
func (s *schedule) reset() {
	s.current = s.nominal
}

// interval returns the current effective interval without advancing state.
func (s *schedule) interval() time.Duration {
	return s.current
}
