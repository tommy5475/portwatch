package daemon

import (
	"errors"
	"testing"
	"time"
)

func TestMetricsInitialState(t *testing.T) {
	m := newMetrics()
	snap := m.snapshot()

	if snap.ScansTotal != 0 {
		t.Errorf("expected ScansTotal 0, got %d", snap.ScansTotal)
	}
	if snap.AlertsTotal != 0 {
		t.Errorf("expected AlertsTotal 0, got %d", snap.AlertsTotal)
	}
	if snap.ChangesTotal != 0 {
		t.Errorf("expected ChangesTotal 0, got %d", snap.ChangesTotal)
	}
	if !snap.LastScanTime.IsZero() {
		t.Errorf("expected zero LastScanTime, got %v", snap.LastScanTime)
	}
}

func TestMetricsRecordScanNoError(t *testing.T) {
	m := newMetrics()
	before := time.Now()
	m.recordScan(nil)
	after := time.Now()

	snap := m.snapshot()
	if snap.ScansTotal != 1 {
		t.Errorf("expected ScansTotal 1, got %d", snap.ScansTotal)
	}
	if snap.LastScanError != "" {
		t.Errorf("expected empty LastScanError, got %q", snap.LastScanError)
	}
	if snap.LastScanTime.Before(before) || snap.LastScanTime.After(after) {
		t.Errorf("LastScanTime %v not in expected range [%v, %v]", snap.LastScanTime, before, after)
	}
}

func TestMetricsRecordScanWithError(t *testing.T) {
	m := newMetrics()
	sentinel := errors.New("scan failed")
	m.recordScan(sentinel)

	snap := m.snapshot()
	if snap.ScansTotal != 1 {
		t.Errorf("expected ScansTotal 1, got %d", snap.ScansTotal)
	}
	if snap.LastScanError != "scan failed" {
		t.Errorf("unexpected LastScanError: %q", snap.LastScanError)
	}
}

func TestMetricsRecordChanges(t *testing.T) {
	m := newMetrics()
	m.recordChanges(3)
	m.recordChanges(0) // should be a no-op
	m.recordChanges(2)

	snap := m.snapshot()
	if snap.ChangesTotal != 5 {
		t.Errorf("expected ChangesTotal 5, got %d", snap.ChangesTotal)
	}
}

func TestMetricsRecordAlert(t *testing.T) {
	m := newMetrics()
	m.recordAlert()
	m.recordAlert()

	snap := m.snapshot()
	if snap.AlertsTotal != 2 {
		t.Errorf("expected AlertsTotal 2, got %d", snap.AlertsTotal)
	}
}

func TestMetricsSnapshotIsIsolated(t *testing.T) {
	m := newMetrics()
	snap1 := m.snapshot()
	m.recordScan(nil)
	snap2 := m.snapshot()

	if snap1.ScansTotal == snap2.ScansTotal {
		t.Error("expected snapshots to differ after recordScan")
	}
}
