package daemon

import "sync"

// queue is a thread-safe, bounded FIFO queue for arbitrary items.
// When the queue is full, the oldest item is evicted to make room.
type queue[T any] struct {
	mu       sync.Mutex
	items    []T
	cap      int
	evicted  int
}

func newQueue[T any](capacity int) *queue[T] {
	if capacity <= 0 {
		capacity = 64
	}
	return &queue[T]{
		items: make([]T, 0, capacity),
		cap:  capacity,
	}
}

// Push adds an item to the back of the queue.
// If the queue is at capacity the oldest item is dropped and the eviction
// counter is incremented.
func (q *queue[T]) Push(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.cap {
		q.items = q.items[1:]
		q.evicted++
	}
	q.items = append(q.items, item)
}

// Pop removes and returns the front item.
// The second return value is false when the queue is empty.
func (q *queue[T]) Pop() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Len returns the current number of items in the queue.
func (q *queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Evicted returns the total number of items dropped due to capacity overflow.
func (q *queue[T]) Evicted() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.evicted
}

// Drain returns all current items and clears the queue.
func (q *queue[T]) Drain() []T {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]T, len(q.items))
	copy(out, q.items)
	q.items = q.items[:0]
	return out
}
