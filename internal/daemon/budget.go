package daemon

import (
	"sync"
	"time"
)

// budget tracks a rolling time-based resource allowance. Each call to Spend
// deducts from the available tokens; tokens replenish at a fixed rate over a
// sliding window. Budget is safe for concurrent use.
type budget struct {
	mu       sync.Mutex
	cap      int
	avail    int
	rate     time.Duration // replenish one token every rate
	lastTick time.Time
}

func newBudget(cap int, rate time.Duration) *budget {
	if cap <= 0 {
		cap = 10
	}
	if rate <= 0 {
		rate = time.Second
	}
	return &budget{
		cap:      cap,
		avail:    cap,
		rate:     rate,
		lastTick: time.Now(),
	}
}

// replenish adds tokens proportional to elapsed time since the last tick.
func (b *budget) replenish(now time.Time) {
	elapsed := now.Sub(b.lastTick)
	if elapsed < b.rate {
		return
	}
	tokens := int(elapsed / b.rate)
	b.avail += tokens
	if b.avail > b.cap {
		b.avail = b.cap
	}
	b.lastTick = b.lastTick.Add(time.Duration(tokens) * b.rate)
}

// Spend attempts to consume n tokens. Returns true if the budget had enough
// tokens; false if the request was denied.
func (b *budget) Spend(n int) bool {
	if n <= 0 {
		return true
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.replenish(time.Now())
	if b.avail < n {
		return false
	}
	b.avail -= n
	return true
}

// Remaining returns the number of tokens currently available.
func (b *budget) Remaining() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.replenish(time.Now())
	return b.avail
}

// Reset refills the budget to its full capacity.
func (b *budget) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.avail = b.cap
	b.lastTick = time.Now()
}
