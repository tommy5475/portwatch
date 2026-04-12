// Package daemon provides the core runtime components for the portwatch daemon.
//
// # Cooldown
//
// cooldown enforces a minimum quiet period between successive activations
// of a named action. It differs from throttle in that it measures time
// from the last *completion* rather than the last *start*, making it
// well-suited for actions with variable execution duration such as
// alert dispatches or webhook deliveries.
//
// Usage:
//
//	cd := newCooldown(30 * time.Second)
//
//	if cd.Ready("webhook") {
//		dispatch()
//		cd.Mark("webhook")
//	}
//
// Cooldown keys are independent; marking one key does not affect others.
// Calling Reset removes a key's record so the next call to Ready returns
// true regardless of when Mark was last called.
package daemon
