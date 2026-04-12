package daemon

import "sync"

// barrier is a reusable synchronisation point that lets N goroutines wait
// until all of them have reached the same checkpoint before any proceeds.
// Once all participants have arrived the barrier is automatically reset so
// it can be used for the next round.
type barrier struct {
	mu      sync.Mutex
	cond    *sync.Cond
	total   int
	arrived int
	gen     uint64 // generation counter – incremented on each release
}

func newBarrier(n int) *barrier {
	if n < 1 {
		n = 1
	}
	b := &barrier{total: n}
	b.cond = sync.NewCond(&b.mu)
	return b
}

// Wait blocks the calling goroutine until all n participants have called
// Wait.  It returns the generation number of the round that just completed.
func (b *barrier) Wait() uint64 {
	b.mu.Lock()
	defer b.mu.Unlock()

	gen := b.gen
	b.arrived++

	if b.arrived == b.total {
		// Last arrival – release everyone and start a new generation.
		b.arrived = 0
		b.gen++
		b.cond.Broadcast()
		return gen
	}

	// Wait until the generation advances.
	for b.gen == gen {
		b.cond.Wait()
	}
	return gen
}

// Size returns the number of participants the barrier was created with.
func (b *barrier) Size() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.total
}

// Arrived returns how many goroutines are currently waiting at the barrier.
func (b *barrier) Arrived() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.arrived
}
