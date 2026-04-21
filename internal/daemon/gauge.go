package daemon

import "sync/atomic"

// gauge is a thread-safe floating-point-like metric that tracks a single
// instantaneous value (e.g. current queue depth, active connections).
// Values are stored as scaled int64 to avoid floating-point atomics.
type gauge struct {
	value atomic.Int64
	min   atomic.Int64
	max   atomic.Int64
	count atomic.Int64
}

const gaugeScale = 1000

func newGauge() *gauge {
	g := &gauge{}
	// initialise min to MaxInt64 so the first Set wins
	g.min.Store(int64(^uint64(0) >> 1))
	return g
}

// Set records v as the current value and updates min/max/count.
func (g *gauge) Set(v float64) {
	scaled := int64(v * gaugeScale)
	g.value.Store(scaled)
	g.count.Add(1)

	for {
		old := g.max.Load()
		if scaled <= old || g.max.CompareAndSwap(old, scaled) {
			break
		}
	}
	for {
		old := g.min.Load()
		if scaled >= old || g.min.CompareAndSwap(old, scaled) {
			break
		}
	}
}

// Get returns the current value.
func (g *gauge) Get() float64 {
	return float64(g.value.Load()) / gaugeScale
}

// Min returns the minimum value recorded since creation or last Reset.
// Returns 0 if Set has never been called.
func (g *gauge) Min() float64 {
	if g.count.Load() == 0 {
		return 0
	}
	return float64(g.min.Load()) / gaugeScale
}

// Max returns the maximum value recorded since creation or last Reset.
func (g *gauge) Max() float64 {
	return float64(g.max.Load()) / gaugeScale
}

// Count returns the number of times Set has been called.
func (g *gauge) Count() int64 {
	return g.count.Load()
}

// Reset clears all state back to initial values.
func (g *gauge) Reset() {
	g.value.Store(0)
	g.count.Store(0)
	g.max.Store(0)
	g.min.Store(int64(^uint64(0) >> 1))
}
