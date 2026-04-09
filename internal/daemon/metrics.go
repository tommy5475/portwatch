package daemon

import (
	"sync/atomic"
	"time"
)

// metrics holds runtime counters for the daemon.
type metrics struct {
	scansTotal    atomic.Int64
	alertsTotal   atomic.Int64
	changesTotal  atomic.Int64
	lastScanTime  atomic.Int64 // Unix nano
	lastScanError atomic.Value // stores last error string or nil
}

func newMetrics() *metrics {
	return &metrics{}
}

// recordScan increments the scan counter and records the timestamp.
// If err is non-nil the error message is stored; otherwise the error
// field is cleared.
func (m *metrics) recordScan(err error) {
	m.scansTotal.Add(1)
	m.lastScanTime.Store(time.Now().UnixNano())
	if err != nil {
		m.lastScanError.Store(err.Error())
	} else {
		m.lastScanError.Store("")
	}
}

// recordChanges increments the changes counter by n.
func (m *metrics) recordChanges(n int) {
	if n > 0 {
		m.changesTotal.Add(int64(n))
	}
}

// recordAlert increments the alert counter.
func (m *metrics) recordAlert() {
	m.alertsTotal.Add(1)
}

// snapshot returns a point-in-time copy of the current metrics.
func (m *metrics) snapshot() MetricsSnapshot {
	snap := MetricsSnapshot{
		ScansTotal:   m.scansTotal.Load(),
		AlertsTotal:  m.alertsTotal.Load(),
		ChangesTotal: m.changesTotal.Load(),
	}
	if ns := m.lastScanTime.Load(); ns != 0 {
		snap.LastScanTime = time.Unix(0, ns)
	}
	if v, ok := m.lastScanError.Load().(string); ok {
		snap.LastScanError = v
	}
	return snap
}

// MetricsSnapshot is an immutable copy of daemon metrics.
type MetricsSnapshot struct {
	ScansTotal   int64
	AlertsTotal  int64
	ChangesTotal int64
	LastScanTime time.Time
	LastScanError string
}
