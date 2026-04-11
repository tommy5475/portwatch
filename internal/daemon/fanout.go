package daemon

import "sync"

// fanout distributes a single input value to multiple registered receivers
// concurrently. Each receiver is a buffered channel; slow receivers that cannot
// accept a value within the buffer are skipped rather than blocking the caller.
type fanout[T any] struct {
	mu       sync.RWMutex
	subs     map[int]chan T
	nextID   int
	bufSize  int
}

func newFanout[T any](bufSize int) *fanout[T] {
	if bufSize < 1 {
		bufSize = 1
	}
	return &fanout[T]{
		subs:    make(map[int]chan T),
		bufSize: bufSize,
	}
}

// subscribe registers a new receiver and returns its channel and a cancel
// function that removes the subscription and closes the channel.
func (f *fanout[T]) subscribe() (<-chan T, func()) {
	f.mu.Lock()
	id := f.nextID
	f.nextID++
	ch := make(chan T, f.bufSize)
	f.subs[id] = ch
	f.mu.Unlock()

	cancel := func() {
		f.mu.Lock()
		delete(f.subs, id)
		close(ch)
		f.mu.Unlock()
	}
	return ch, cancel
}

// publish sends v to all current subscribers. Subscribers whose buffer is full
// are skipped so that one slow consumer cannot stall others.
func (f *fanout[T]) publish(v T) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, ch := range f.subs {
		select {
		case ch <- v:
		default:
		}
	}
}

// len returns the current number of subscribers.
func (f *fanout[T]) len() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.subs)
}
