package daemon

import (
	"context"
	"sync/atomic"
	"time"
)

// watcherState tracks whether the watcher considers the system healthy.
type watcherState struct {
	consecutiveFailures atomic.Int64
	lastSuccess         atomic.Value // stores time.Time
	threshold           int64
}

// watcher monitors the scan pipeline for consecutive failures and
// transitions the daemon into a degraded state when the failure
// threshold is exceeded.
type watcher struct {
	state     *watcherState
	threshold int64
	onDegrade func()
	onRecover func()
}

func newWatcher(threshold int64, onDegrade, onRecover func()) *watcher {
	if threshold <= 0 {
		threshold = 3
	}
	if onDegrade == nil {
		onDegrade = func() {}
	}
	if onRecover == nil {
		onRecover = func() {}
	}
	w := &watcher{
		state:     &watcherState{threshold: threshold},
		threshold: threshold,
		onDegrade: onDegrade,
		onRecover: onRecover,
	}
	w.state.lastSuccess.Store(time.Time{})
	return w
}

// recordSuccess resets the failure counter and fires onRecover if the
// system was previously degraded.
func (w *watcher) recordSuccess() {
	prev := w.state.consecutiveFailures.Swap(0)
	w.state.lastSuccess.Store(time.Now())
	if prev >= w.threshold {
		w.onRecover()
	}
}

// recordFailure increments the failure counter and fires onDegrade when
// the threshold is first crossed.
func (w *watcher) recordFailure() {
	next := w.state.consecutiveFailures.Add(1)
	if next == w.threshold {
		w.onDegrade()
	}
}

// isDegraded reports whether the failure threshold has been reached.
func (w *watcher) isDegraded() bool {
	return w.state.consecutiveFailures.Load() >= w.threshold
}

// lastSuccessTime returns the time of the most recent successful scan.
func (w *watcher) lastSuccessTime() time.Time {
	return w.state.lastSuccess.Load().(time.Time)
}

// runUntil blocks until ctx is cancelled, calling check every interval.
func (w *watcher) runUntil(ctx context.Context, interval time.Duration, check func() error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := check(); err != nil {
				w.recordFailure()
			} else {
				w.recordSuccess()
			}
		}
	}
}
