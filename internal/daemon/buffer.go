package daemon

import "sync"

// buffer is a bounded, thread-safe ring buffer that stores the most recent N
// items. When the buffer is full the oldest entry is overwritten.
type buffer[T any] struct {
	mu    sync.Mutex
	items []T
	head  int
	count int
	cap   int
}

// newBuffer creates a ring buffer with the given capacity.
// If cap is less than 1 it defaults to 16.
func newBuffer[T any](cap int) *buffer[T] {
	if cap < 1 {
		cap = 16
	}
	return &buffer[T]{
		items: make([]T, cap),
		cap:   cap,
	}
}

// Push adds an item to the buffer, overwriting the oldest entry when full.
func (b *buffer[T]) Push(item T) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.items[b.head] = item
	b.head = (b.head + 1) % b.cap
	if b.count < b.cap {
		b.count++
	}
}

// Snapshot returns a copy of all stored items in insertion order (oldest first).
func (b *buffer[T]) Snapshot() []T {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.count == 0 {
		return nil
	}
	out := make([]T, b.count)
	start := (b.head - b.count + b.cap) % b.cap
	for i := 0; i < b.count; i++ {
		out[i] = b.items[(start+i)%b.cap]
	}
	return out
}

// Len returns the number of items currently stored.
func (b *buffer[T]) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.count
}

// Reset clears all items from the buffer.
func (b *buffer[T]) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.head = 0
	b.count = 0
}
