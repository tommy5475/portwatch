// Package daemon provides the core runtime primitives for the portwatch daemon.
//
// ticker.go implements a resilient scan ticker that drives the main polling loop.
// It wraps time.Ticker with jitter support and cooperative cancellation via
// context, ensuring clean shutdown without goroutine leaks.
//
// Usage:
//
//	t := newTicker(ctx, interval, jitter)
//	runLoop(ctx, t, func(ctx context.Context) error {
//		return scanPorts(ctx)
//	})
//
// The jitter fraction (0.0–1.0) adds a random offset to each tick interval,
// spreading load when multiple portwatch instances run in parallel.
package daemon
