package daemon

import (
	"sync"
	"time"
)

// relay forwards events from one channel to many subscribers with optional
// filtering and a configurable drop policy for slow consumers.
type relay[T any] struct {
	mu      sync.RWMutex
	subs    map[int]chan T
	next    int
	bufSize int
	filter  func(T) bool
}

func newRelay[T any](bufSize int, filter func(T) bool) *relay[T] {
	if bufSize < 1 {
		bufSize = 16
	}
	if filter == nil {
		filter = func(T) bool { return true }
	}
	return &relay[T]{
		subs:    make(map[int]chan T),
		bufSize: bufSize,
		filter:  filter,
	}
}

// subscribe registers a new subscriber and returns its channel and a cancel func.
func (r *relay[T]) subscribe() (<-chan T, func()) {
	r.mu.Lock()
	id := r.next
	r.next++
	ch := make(chan T, r.bufSize)
	r.subs[id] = ch
	r.mu.Unlock()
	return ch, func() {
		r.mu.Lock()
		delete(r.subs, id)
		close(ch)
		r.mu.Unlock()
	}
}

// send delivers v to all subscribers that pass the filter.
// Slow consumers are skipped after a short non-blocking attempt.
func (r *relay[T]) send(v T) {
	if !r.filter(v) {
		return
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ch := range r.subs {
		select {
		case ch <- v:
		case <-time.After(time.Millisecond):
		}
	}
}

// len returns the current subscriber count.
func (r *relay[T]) len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.subs)
}
