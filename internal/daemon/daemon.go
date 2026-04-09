// Package daemon wires together the scanner, monitor, filter, notifier,
// and reporter into a long-running portwatch process.
package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Daemon orchestrates all portwatch subsystems.
type Daemon struct {
	cfg      *config.Config
	mon      *monitor.Monitor
	not      *notifier.Notifier
	rep      *reporter.Reporter
	filter   *filter.Filter
}

// New constructs a Daemon from the supplied configuration.
func New(cfg *config.Config) (*Daemon, error) {
	sc := scanner.New(time.Duration(cfg.TimeoutMS) * time.Millisecond)

	st, err := state.New(cfg.StateFile)
	if err != nil {
		return nil, err
	}

	f, err := filter.New(cfg.Filter.Ports, cfg.Filter.Protocols, cfg.Filter.DefaultPermit)
	if err != nil {
		return nil, err
	}

	mon := monitor.New(sc, st, cfg)
	not := notifier.New(cfg)
	rep := reporter.New(cfg.ReportFormat)

	return &Daemon{
		cfg:    cfg,
		mon:    mon,
		not:    not,
		rep:    rep,
		filter: f,
	}, nil
}

// Run starts the daemon and blocks until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) error {
	if err := d.mon.Start(ctx); err != nil {
		return err
	}

	ch := d.mon.Changes()
	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch: shutting down")
			return ctx.Err()
		case diff, ok := <-ch:
			if !ok {
				return nil
			}
			filtered := d.filter.Apply(diff)
			if len(filtered.Opened)+len(filtered.Closed) == 0 {
				continue
			}
			output := d.rep.Report(filtered)
			if err := d.not.Notify(ctx, output, filtered); err != nil {
				log.Printf("portwatch: notify error: %v", err)
			}
		}
	}
}
