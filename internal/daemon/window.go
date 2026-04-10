package daemon

import (
	"sync"
	"time"
)

// window tracks event counts within a sliding time window.
// It is safe for concurrent use.
type window struct {
	mu       sync.Mutex
	buckets  []int
	size     int
	interval time.Duration
	last     time.Time
}

// newWindow creates a sliding window divided into `size` buckets,
// each covering `interval` duration. Total window span = size * interval.
// Falls back to safe defaults when arguments are invalid.
func newWindow(size int, interval time.Duration) *window {
	if size <= 0 {
		size = 10
	}
	if interval <= 0 {
		interval = time.Second
	}
	return &window{
		buckets:  make([]int, size),
		size:     size,
		interval: interval,
		last:     time.Now(),
	}
}

// Record increments the count for the current time bucket,
// advancing (and clearing) stale buckets as needed.
func (w *window) Record() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(time.Now())
	w.buckets[0]++
}

// Count returns the total number of events recorded across all live buckets.
func (w *window) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(time.Now())
	total := 0
	for _, v := range w.buckets {
		total += v
	}
	return total
}

// Reset clears all buckets.
func (w *window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := range w.buckets {
		w.buckets[i] = 0
	}
	w.last = time.Now()
}

// advance rotates buckets based on elapsed time since last call.
// Must be called with w.mu held.
func (w *window) advance(now time.Time) {
	elapsed := now.Sub(w.last)
	steps := int(elapsed / w.interval)
	if steps <= 0 {
		return
	}
	w.last = w.last.Add(time.Duration(steps) * w.interval)
	if steps >= w.size {
		for i := range w.buckets {
			w.buckets[i] = 0
		}
		return
	}
	// Rotate: shift existing buckets right by `steps`, zero the new front slots.
	copy(w.buckets[steps:], w.buckets[:w.size-steps])
	for i := 0; i < steps; i++ {
		w.buckets[i] = 0
	}
}
