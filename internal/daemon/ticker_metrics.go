package daemon

import (
	"sync/atomic"
	"time"
)

// tickerMetrics tracks runtime statistics for the scan loop ticker.
type tickerMetrics struct {
	tickCount    atomic.Int64
	skipCount    atomic.Int64
	lastTickedAt atomic.Int64 // UnixNano
	totalLag     atomic.Int64 // nanoseconds
}

func newTickerMetrics() *tickerMetrics {
	return &tickerMetrics{}
}

// recordTick notes a tick occurred at the given time, recording lag vs expected.
func (m *tickerMetrics) recordTick(at time.Time, lag time.Duration) {
	m.tickCount.Add(1)
	m.lastTickedAt.Store(at.UnixNano())
	if lag > 0 {
		m.totalLag.Add(lag.Nanoseconds())
	}
}

// recordSkip notes a tick was skipped (e.g. previous scan still running).
func (m *tickerMetrics) recordSkip() {
	m.skipCount.Add(1)
}

// TickCount returns total ticks processed.
func (m *tickerMetrics) TickCount() int64 {
	return m.tickCount.Load()
}

// SkipCount returns total ticks skipped.
func (m *tickerMetrics) SkipCount() int64 {
	return m.skipCount.Load()
}

// LastTickedAt returns the time of the most recent tick, or zero if none.
func (m *tickerMetrics) LastTickedAt() time.Time {
	ns := m.lastTickedAt.Load()
	if ns == 0 {
		return time.Time{}
	}
	return time.Unix(0, ns)
}

// AverageLag returns the mean lag per tick, or zero if no ticks recorded.
func (m *tickerMetrics) AverageLag() time.Duration {
	count := m.tickCount.Load()
	if count == 0 {
		return 0
	}
	return time.Duration(m.totalLag.Load() / count)
}

// reset clears all counters.
func (m *tickerMetrics) reset() {
	m.tickCount.Store(0)
	m.skipCount.Store(0)
	m.lastTickedAt.Store(0)
	m.totalLag.Store(0)
}
