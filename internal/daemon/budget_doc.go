// Package daemon contains the internal runtime components of portwatch.
//
// # Budget
//
// budget provides a token-bucket style resource allowance that replenishes
// over time. It is used to cap expensive operations (e.g. alert dispatches or
// webhook calls) to a configurable rate without hard-blocking callers.
//
// Usage:
//
//	b := newBudget(10, time.Second) // 10 tokens, one replenished per second
//	if b.Spend(1) {
//	    // perform the operation
//	}
//
// Tokens accumulate up to the configured capacity; they are never over-filled.
// All methods are safe for concurrent use.
package daemon
