// Package daemon provides the core runtime components of portwatch.
//
// # Queue
//
// queue[T] is a generic, thread-safe bounded FIFO queue.
//
// It is designed for internal use where producers can outpace consumers for
// short bursts. When the queue reaches its configured capacity the oldest
// entry is silently evicted and the eviction counter is incremented so callers
// can detect back-pressure without blocking.
//
// Typical usage:
//
//	q := newQueue[string](256)
//	q.Push("event-a")
//	q.Push("event-b")
//
//	for {
//		item, ok := q.Pop()
//		if !ok {
//			break
//		}
//		// process item
//	}
//
// Evictions can be monitored via q.Evicted().
package daemon
