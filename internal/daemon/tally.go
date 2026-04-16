package daemon

import "sync"

// tally tracks named integer counters with optional ceiling.
type tally struct {
	mu      sync.Mutex
	counts  map[string]int64
	ceiling int64
}

func newTally(ceiling int64) *tally {
	if ceiling <= 0 {
		ceiling = 1<<63 - 1
	}
	return &tally{
		counts:  make(map[string]int64),
		ceiling: ceiling,
	}
}

// Inc increments the named counter by 1, capped at ceiling.
func (t *tally) Inc(key string) int64 {
	return t.Add(key, 1)
}

// Add adds delta to the named counter, capped at ceiling.
func (t *tally) Add(key string, delta int64) int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	v := t.counts[key] + delta
	if v > t.ceiling {
		v = t.ceiling
	}
	if v < 0 {
		v = 0
	}
	t.counts[key] = v
	return v
}

// Get returns the current value for key.
func (t *tally) Get(key string) int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.counts[key]
}

// Reset sets the named counter to zero.
func (t *tally) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.counts, key)
}

// Snapshot returns a copy of all counters.
func (t *tally) Snapshot() map[string]int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[string]int64, len(t.counts))
	for k, v := range t.counts {
		out[k] = v
	}
	return out
}

// Len returns the number of tracked keys.
func (t *tally) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.counts)
}
