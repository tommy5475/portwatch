// Package daemon — relay
//
// relay is a generic, filter-aware fan-out primitive that forwards values of
// type T from a single producer to an arbitrary number of subscribers.
//
// Subscribers register via subscribe(), which returns a buffered receive
// channel and a cancel function. Calling cancel unregisters the subscriber
// and closes its channel.
//
// An optional filter predicate can be supplied at construction time; events
// that do not pass the predicate are silently dropped before any subscriber
// sees them.
//
// Slow consumers are handled gracefully: send performs a non-blocking
// delivery attempt and skips any subscriber whose buffer is full, preventing
// a lagging subscriber from stalling the producer.
//
// relay is safe for concurrent use.
package daemon
