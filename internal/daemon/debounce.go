package daemon

import (
	"sync"
	"time"
)

// debounce delays execution of a function until after a quiet period has
// elapsed since the last invocation. Repeated calls within the wait duration
// reset the timer, ensuring the function runs only once per burst of events.
type debounce struct {
	mu    sync.Mutex
	wait  time.Duration
	timer *time.Timer
	fn    func()
}

// newDebounce creates a debounce with the given quiet-period duration.
// If wait is <= 0 it defaults to 500ms.
func newDebounce(wait time.Duration, fn func()) *debounce {
	if wait <= 0 {
		wait = 500 * time.Millisecond
	}
	return &debounce{
		wait: wait,
		fn:   fn,
	}
}

// Trigger schedules fn to be called after the quiet period. If Trigger is
// called again before the period elapses the timer resets.
func (d *debounce) Trigger() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.wait, d.fn)
}

// Flush cancels any pending timer and calls fn immediately. It is a no-op
// if no call is pending.
func (d *debounce) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil && d.timer.Stop() {
		go d.fn()
		d.timer = nil
	}
}

// Stop cancels any pending timer without invoking fn.
func (d *debounce) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
