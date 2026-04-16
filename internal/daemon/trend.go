package daemon

import (
	"sync"
	"time"
)

// trend tracks a rolling rate of events over a sliding window,
// allowing callers to observe whether activity is increasing or decreasing.
type trend struct {
	mu       sync.Mutex
	buckets  []int64
	size     int
	interval time.Duration
	start    time.Time
}

func newTrend(buckets int, interval time.Duration) *trend {
	if buckets < 2 {
		buckets = 10
	}
	if interval <= 0 {
		interval = time.Minute
	}
	return &trend{
		buckets:  make([]int64, buckets),
		size:     buckets,
		interval: interval,
		start:    time.Now(),
	}
}

// Record adds n events to the current bucket.
func (t *trend) record(n int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	idx := t.bucketIndex(time.Now())
	t.buckets[idx] += n
}

// Rate returns the total count across all buckets.
func (t *trend) rate() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	var total int64
	for _, v := range t.buckets {
		total += v
	}
	return total
}

// Rising returns true if the second half of buckets exceeds the first half.
func (t *trend) rising() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	mid := t.size / 2
	var first, second int64
	for i := 0; i < mid; i++ {
		first += t.buckets[i]
	}
	for i := mid; i < t.size; i++ {
		second += t.buckets[i]
	}
	return second > first
}

// Reset clears all buckets.
func (t *trend) reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i := range t.buckets {
		t.buckets[i] = 0
	}
	t.start = time.Now()
}

func (t *trend) bucketIndex(now time.Time) int {
	elapsed := now.Sub(t.start)
	idx := int(elapsed/t.interval) % t.size
	if idx < 0 {
		idx = 0
	}
	return idx
}
