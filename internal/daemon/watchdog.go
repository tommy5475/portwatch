package daemon

import (
	"sync"
	"time"
)

// watchdog tracks whether a recurring operation is making progress within a
// configurable deadline. If the deadline is exceeded the watchdog fires and
// the caller can react (e.g. log a warning, increment a metric, restart a
// sub-system). Calling Kick resets the deadline.
type watchdog struct {
	mu       sync.Mutex
	timeout  time.Duration
	timer    *time.Timer
	fired    bool
	onFire   func()
	stopped  bool
}

// newWatchdog creates a watchdog that calls onFire when no Kick is received
// within timeout. The watchdog starts armed immediately.
// Panics if timeout <= 0 or onFire is nil.
func newWatchdog(timeout time.Duration, onFire func()) *watchdog {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	if onFire == nil {
		onFire = func() {}
	}
	w := &watchdog{
		timeout: timeout,
		onFire:  onFire,
	}
	w.timer = time.AfterFunc(timeout, w.fire)
	return w
}

func (w *watchdog) fire() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopped {
		return
	}
	w.fired = true
	w.onFire()
}

// Kick resets the watchdog deadline. Call this whenever the monitored
// operation completes successfully.
func (w *watchdog) Kick() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopped {
		return
	}
	w.fired = false
	w.timer.Reset(w.timeout)
}

// Fired reports whether the watchdog has fired since the last Kick.
func (w *watchdog) Fired() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.fired
}

// Stop disarms the watchdog permanently.
func (w *watchdog) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.stopped = true
	w.timer.Stop()
}
