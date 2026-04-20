package daemon

import (
	"sync"
	"time"
)

// dial tracks outbound connection attempt latency and success/failure counts
// for a named target. It is safe for concurrent use.
type dial struct {
	mu       sync.Mutex
	target   string
	attempts int64
	successes int64
	failures  int64
	lastAt   time.Time
	lastOK   bool
	totalLat time.Duration
}

func newDial(target string) *dial {
	if target == "" {
		target = "unknown"
	}
	return &dial{target: target}
}

// record registers one connection attempt with its outcome and round-trip duration.
func (d *dial) record(ok bool, lat time.Duration) {
	if lat < 0 {
		lat = 0
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.attempts++
	d.lastAt = time.Now()
	d.lastOK = ok
	d.totalLat += lat
	if ok {
		d.successes++
	} else {
		d.failures++
	}
}

// attempts returns total number of recorded attempts.
func (d *dial) count() int64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.attempts
}

// successRate returns the fraction of successful attempts in [0,1].
// Returns 0 if no attempts have been recorded.
func (d *dial) successRate() float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.attempts == 0 {
		return 0
	}
	return float64(d.successes) / float64(d.attempts)
}

// avgLatency returns the mean round-trip duration across all recorded attempts.
func (d *dial) avgLatency() time.Duration {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.attempts == 0 {
		return 0
	}
	return d.totalLat / time.Duration(d.attempts)
}

// lastSuccess reports whether the most recent attempt succeeded.
func (d *dial) lastSuccess() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.lastOK
}

// reset clears all recorded state.
func (d *dial) reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.attempts = 0
	d.successes = 0
	d.failures = 0
	d.totalLat = 0
	d.lastAt = time.Time{}
	d.lastOK = false
}
