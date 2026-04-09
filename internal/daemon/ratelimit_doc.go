// Package daemon contains the runtime orchestration for portwatch.
//
// # Rate Limiting
//
// The rateLimiter type (ratelimit.go) provides a simple token-bucket
// implementation used to throttle outbound alert notifications.
//
// When a scan cycle detects port changes the daemon passes each change
// through the rate limiter before forwarding it to the notifier. This
// prevents a sudden burst of changes (e.g. a host restart) from
// overwhelming a webhook endpoint or filling log output.
//
// Usage:
//
//	rl := newRateLimiter(10, time.Minute)
//	if rl.Allow() {
//		// send alert
//	}
//
// Tokens are replenished proportionally as time passes, so sustained
// low-frequency changes are always permitted while short-lived spikes
// are smoothed out.
package daemon
