package daemon

import "sync/atomic"

// counter is a thread-safe monotonic counter used to track cumulative
// occurrences of events such as scans completed, alerts fired, or errors
// encountered during a daemon run.
type counter struct {
	value atomic.Int64
}

// newCounter returns a new counter initialised to zero.
func newCounter() *counter {
	return &counter{}
}

// Inc increments the counter by one and returns the new value.
func (c *counter) Inc() int64 {
	return c.value.Add(1)
}

// Add increments the counter by delta and returns the new value.
// delta may be negative to decrement, though the counter does not
// enforce non-negativity.
func (c *counter) Add(delta int64) int64 {
	return c.value.Add(delta)
}

// Get returns the current counter value without modifying it.
func (c *counter) Get() int64 {
	return c.value.Load()
}

// Reset sets the counter back to zero and returns the value that was
// held immediately before the reset.
func (c *counter) Reset() int64 {
	return c.value.Swap(0)
}
