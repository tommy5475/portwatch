package daemon

import (
	"sync"
	"time"
)

// circuitState represents the current state of the circuit breaker.
type circuitState int

const (
	circuitClosed   circuitState = iota // normal operation
	circuitOpen                         // blocking calls
	circuitHalfOpen                     // probing for recovery
)

// circuit is a simple circuit breaker that opens after a threshold of
// consecutive failures and resets after a cooldown period.
type circuit struct {
	mu           sync.Mutex
	state        circuitState
	failures      int
	threshold    int
	cooldown     time.Duration
	openedAt     time.Time
	totalTrips   int
}

func newCircuit(threshold int, cooldown time.Duration) *circuit {
	if threshold <= 0 {
		threshold = 3
	}
	if cooldown <= 0 {
		cooldown = 30 * time.Second
	}
	return &circuit{
		state:     circuitClosed,
		threshold: threshold,
		cooldown:  cooldown,
	}
}

// allow reports whether the caller should proceed. It transitions
// an open circuit to half-open once the cooldown has elapsed.
func (c *circuit) allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch c.state {
	case circuitClosed:
		return true
	case circuitOpen:
		if time.Since(c.openedAt) >= c.cooldown {
			c.state = circuitHalfOpen
			return true
		}
		return false
	case circuitHalfOpen:
		return true
	}
	return false
}

// recordSuccess resets the circuit to closed.
func (c *circuit) recordSuccess() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures = 0
	c.state = circuitClosed
}

// recordFailure increments the failure counter and opens the circuit
// when the threshold is reached.
func (c *circuit) recordFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures++
	if c.failures >= c.threshold && c.state != circuitOpen {
		c.state = circuitOpen
		c.openedAt = time.Now()
		c.totalTrips++
	}
}

// stats returns a snapshot of circuit breaker counters.
func (c *circuit) stats() (state circuitState, failures, trips int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state, c.failures, c.totalTrips
}
