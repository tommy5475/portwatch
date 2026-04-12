package daemon

import "sync"

// ringbuf is a fixed-capacity circular buffer that overwrites the oldest
// entry when full. It is safe for concurrent use.
type ringbuf[T any] struct {
	mu   sync.Mutex
	buf  []T
	head int // index of next write
	len  int // number of valid entries
	cap  int
}

const defaultRingbufCap = 64

// newRingbuf returns a ringbuf with the given capacity.
// If cap <= 0 the default capacity is used.
func newRingbuf[T any](cap int) *ringbuf[T] {
	if cap <= 0 {
		cap = defaultRingbufCap
	}
	return &ringbuf[T]{buf: make([]T, cap), cap: cap}
}

// Push appends v to the buffer, overwriting the oldest entry when full.
func (r *ringbuf[T]) Push(v T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buf[r.head] = v
	r.head = (r.head + 1) % r.cap
	if r.len < r.cap {
		r.len++
	}
}

// Snapshot returns a copy of all valid entries in insertion order
// (oldest first).
func (r *ringbuf[T]) Snapshot() []T {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.len == 0 {
		return nil
	}
	out := make([]T, r.len)
	start := (r.head - r.len + r.cap) % r.cap
	for i := 0; i < r.len; i++ {
		out[i] = r.buf[(start+i)%r.cap]
	}
	return out
}

// Len returns the number of entries currently stored.
func (r *ringbuf[T]) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.len
}

// Reset clears all entries.
func (r *ringbuf[T]) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.head = 0
	r.len = 0
}
