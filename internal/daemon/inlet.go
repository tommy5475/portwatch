package daemon

import (
	"sync"
	"time"
)

// inlet is a bounded, non-blocking ingress point that accepts values of any
// type and makes them available for downstream consumption. Values that arrive
// when the inlet is full are silently dropped so that producers never block.
//
// inlet is safe for concurrent use.
type inlet[T any] struct {
	mu       sync.Mutex
	ch       chan T
	dropped  uint64
	accepted uint64
	created  time.Time
}

func newInlet[T any](capacity int) *inlet[T] {
	if capacity < 1 {
		capacity = 16
	}
	return &inlet[T]{
		ch:      make(chan T, capacity),
		created: time.Now(),
	}
}

// Send attempts to enqueue v. Returns true if v was accepted, false if dropped.
func (in *inlet[T]) Send(v T) bool {
	select {
	case in.ch <- v:
		in.mu.Lock()
		in.accepted++
		in.mu.Unlock()
		return true
	default:
		in.mu.Lock()
		in.dropped++
		in.mu.Unlock()
		return false
	}
}

// Out returns the read-only channel that consumers should read from.
func (in *inlet[T]) Out() <-chan T {
	return in.ch
}

// Dropped returns the total number of values dropped due to a full buffer.
func (in *inlet[T]) Dropped() uint64 {
	in.mu.Lock()
	defer in.mu.Unlock()
	return in.dropped
}

// Accepted returns the total number of values successfully enqueued.
func (in *inlet[T]) Accepted() uint64 {
	in.mu.Lock()
	defer in.mu.Unlock()
	return in.accepted
}

// Len returns the number of values currently buffered.
func (in *inlet[T]) Len() int {
	return len(in.ch)
}

// Age returns the duration since the inlet was created.
func (in *inlet[T]) Age() time.Duration {
	return time.Since(in.created)
}

// Close closes the underlying channel. Callers must not call Send after Close.
func (in *inlet[T]) Close() {
	close(in.ch)
}
