package daemon

import (
	"context"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/state"
)

// pipeline wires together a single scan-diff-notify cycle.
// It is intentionally stateless between calls so that the ticker
// can invoke it repeatedly without shared mutable fields.
type pipeline struct {
	mon      *monitor.Monitor
	filter   *filter.Filter
	store    *state.Store
	notify   *notifier.Notifier
	metrics  *metrics
	snapshot *snapshot
}

func newPipeline(
	mon *monitor.Monitor,
	f *filter.Filter,
	store *state.Store,
	n *notifier.Notifier,
	m *metrics,
	snap *snapshot,
) *pipeline {
	return &pipeline{
		mon:      mon,
		filter:   f,
		store:    store,
		notify:   n,
		metrics:  m,
		snapshot: snap,
	}
}

// run executes one full scan cycle: scan → filter → diff → notify → persist.
// It returns a non-nil error only when the scan itself fails; diff/notify
// errors are recorded in metrics but do not abort the cycle.
func (p *pipeline) run(ctx context.Context) error {
	start := time.Now()

	current, err := p.mon.Scan(ctx)
	if err != nil {
		p.metrics.recordScan(time.Since(start), err)
		return fmt.Errorf("pipeline scan: %w", err)
	}

	filtered := p.filter.Apply(current)
	p.snapshot.update(filtered)
	p.metrics.recordScan(time.Since(start), nil)

	previous, _ := p.store.Load()
	diff := state.Diff(previous, filtered)

	if len(diff.Added)+len(diff.Removed) == 0 {
		return nil
	}

	p.metrics.recordChanges(diff)

	if notifyErr := p.notify.Notify(ctx, diff); notifyErr != nil {
		p.metrics.recordAlert(notifyErr)
	} else {
		p.metrics.recordAlert(nil)
	}

	if saveErr := p.store.Save(filtered); saveErr != nil {
		return fmt.Errorf("pipeline persist: %w", saveErr)
	}

	return nil
}
