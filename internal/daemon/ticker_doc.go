// Package daemon provides the core runtime components for the portwatch daemon.
//
// # Ticker
//
// The ticker module drives the periodic scan loop. It wraps time.Ticker to
// allow deterministic testing via a swappable clock interface.
//
// A [newTicker] call returns a ticker that fires at the configured interval.
// Callers pass a context; when the context is cancelled the ticker stops
// automatically and the run loop exits cleanly.
//
// The [runLoop] function ties the ticker to the pipeline: on every tick it
// invokes the provided scan function and forwards any error to the watcher
// for health tracking.
//
// Usage:
//
//	t := newTicker(interval)
//	defer t.stop()
//	runLoop(ctx, t, interval, func(ctx context.Context) error {
//		return pipeline.Run(ctx)
//	})
package daemon
