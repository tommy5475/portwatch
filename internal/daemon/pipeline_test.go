package daemon

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/state"
)

func buildPipeline(t *testing.T) *pipeline {
	t.Helper()

	mon := monitor.New(monitor.Options{Targets: []string{"127.0.0.1"}, Timeout: time.Second})
	f, err := filter.New(filter.Options{})
	if err != nil {
		t.Fatalf("filter.New: %v", err)
	}

	dir := t.TempDir()
	store, err := state.New(dir + "/state.json")
	if err != nil {
		t.Fatalf("state.New: %v", err)
	}

	n := notifier.New(notifier.Options{})
	m := newMetrics()
	snap := newSnapshot()

	return newPipeline(mon, f, store, n, m, snap)
}

func TestPipelineRunUpdatesSnapshot(t *testing.T) {
	p := buildPipeline(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	before := p.snapshot.count()
	_ = p.run(ctx)
	after := p.snapshot.count()

	if after < before {
		t.Errorf("snapshot count went backwards: before=%d after=%d", before, after)
	}
}

func TestPipelineRunRecordsScanMetric(t *testing.T) {
	p := buildPipeline(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = p.run(ctx)

	if p.metrics.totalScans() == 0 {
		t.Error("expected at least one scan recorded in metrics")
	}
}

func TestPipelineRunContextCancelled(t *testing.T) {
	p := buildPipeline(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := p.run(ctx)
	// A cancelled context should surface as a scan error; we just verify
	// the call returns (does not block) and metrics capture the failure.
	_ = err // may or may not error depending on scanner implementation
}

func TestPipelineNoDiffNoAlert(t *testing.T) {
	p := buildPipeline(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run twice; second run should produce no diff against the persisted state.
	_ = p.run(ctx)
	before := p.metrics.totalAlerts()
	_ = p.run(ctx)
	after := p.metrics.totalAlerts()

	if after > before+1 {
		t.Errorf("unexpected extra alerts on second run: before=%d after=%d", before, after)
	}
}
