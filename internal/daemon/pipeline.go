package daemon

import (
	"context"
	"fmt"
	"log"

	"portwatch/internal/alert"
	"portwatch/internal/monitor"
	"portwatch/internal/state"
)

// pipeline wires together one scan cycle: scan → diff → throttle → alert.
type pipeline struct {
	mon      *monitor.Monitor
	store    *state.Store
	alerter  *alert.Alerter
	throttle *throttle
	snap     *snapshot
	metrics  *metrics
}

func newPipeline(
	mon *monitor.Monitor,
	store *state.Store,
	alerter *alert.Alerter,
	th *throttle,
	snap *snapshot,
	m *metrics,
) *pipeline {
	return &pipeline{
		mon:     mon,
		store:   store,
		alerter: alerter,
		throttle: th,
		snap:    snap,
		metrics: m,
	}
}

// Run executes one full scan cycle. It is intended to be called on each
// ticker tick inside the daemon run-loop.
func (p *pipeline) Run(ctx context.Context) error {
	ports, err := p.mon.Scan(ctx)
	p.metrics.recordScan(err)
	if err != nil {
		return fmt.Errorf("scan: %w", err)
	}

	p.snap.Update(ports)

	prev, err := p.store.Load()
	if err != nil {
		log.Printf("[pipeline] state load: %v — treating as empty", err)
		prev = state.New()
	}

	curr := state.New()
	for _, port := range ports {
		curr.Add(port)
	}

	diff := prev.Diff(curr)
	if diff.Empty() {
		return p.store.Save(curr)
	}

	p.metrics.recordChanges(diff)

	for _, change := range diff.Changes() {
		key := fmt.Sprintf("%s:%d:%s", change.Protocol, change.Port, change.Kind)
		if !p.throttle.Allow(key) {
			log.Printf("[pipeline] throttled alert for %s", key)
			continue
		}
		if alertErr := p.alerter.Send(ctx, change); alertErr != nil {
			log.Printf("[pipeline] alert send: %v", alertErr)
			p.metrics.recordAlert(alertErr)
		} else {
			p.metrics.recordAlert(nil)
		}
	}

	return p.store.Save(curr)
}
