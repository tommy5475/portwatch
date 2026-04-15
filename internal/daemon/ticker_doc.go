// Package daemon provides the core runtime primitives for the portwatch daemon.
//
// # Ticker
//
// ticker is a thin wrapper around [time.Ticker] that exposes a channel-based
// interface consumed by [runLoop]. It exists primarily to allow deterministic
// unit testing via [newFakeTicker], which lets tests drive ticks manually
// without relying on wall-clock timing.
//
// # runLoop
//
// runLoop drives the main scan cycle. It blocks until the provided
// [context.Context] is cancelled, calling fn on every tick emitted by the
// supplied ticker. The function is called synchronously; the loop will not
// tick again until fn returns, ensuring that slow scans do not overlap.
//
// Typical usage:
//
//	tk := newTicker(cfg.Interval)
//	defer tk.stop()
//	runLoop(ctx, tk, func() {
//		_ = pipeline.run(ctx)
//	})
package daemon
