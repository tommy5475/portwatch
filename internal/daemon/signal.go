package daemon

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// RunWithSignals starts the daemon and returns when SIGINT or SIGTERM is
// received, or when the parent context is cancelled.
//
// It is a convenience wrapper intended for use in main().
func RunWithSignals(parent context.Context, d *Daemon) error {
	ctx, stop := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("portwatch: started — press Ctrl+C to stop")
	err := d.Run(ctx)

	// A cancelled context due to a signal is not an application error.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	return err
}
