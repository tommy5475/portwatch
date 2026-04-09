package daemon

import (
	"context"
	"time"
)

// ticker wraps time.Ticker to provide a testable interval-based tick source.
type ticker struct {
	t        *time.Ticker
	interval time.Duration
}

// newTicker creates a ticker that fires at the given interval.
func newTicker(interval time.Duration) *ticker {
	return &ticker{
		t:        time.NewTicker(interval),
		interval: interval,
	}
}

// C returns the underlying channel that delivers ticks.
func (tk *ticker) C() <-chan time.Time {
	return tk.t.C
}

// Stop halts the ticker, releasing associated resources.
func (tk *ticker) Stop() {
	tk.t.Stop()
}

// tickSource is the interface used by the daemon loop so tests can inject a
// fake implementation without real wall-clock delays.
type tickSource interface {
	C() <-chan time.Time
	Stop()
}

// runLoop executes fn on every tick from src until ctx is cancelled.
// It returns the context error that caused it to stop.
func runLoop(ctx context.Context, src tickSource, fn func()) error {
	defer src.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-src.C():
			fn()
		}
	}
}
