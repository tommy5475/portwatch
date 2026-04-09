package daemon_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
)

func minimalConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.Default()
	cfg.StateFile = t.TempDir() + "/state.json"
	return cfg
}

func TestNew(t *testing.T) {
	cfg := minimalConfig(t)
	d, err := daemon.New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if d == nil {
		t.Fatal("New() returned nil daemon")
	}
}

func TestNewInvalidFilter(t *testing.T) {
	cfg := minimalConfig(t)
	cfg.Filter.Ports = []string{"not-a-port"}
	_, err := daemon.New(cfg)
	if err == nil {
		t.Fatal("expected error for invalid filter port range, got nil")
	}
}

func TestRunCancelImmediately(t *testing.T) {
	cfg := minimalConfig(t)
	cfg.Interval = 60 // long interval so no scan fires
	d, err := daemon.New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = d.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Fatalf("Run() unexpected error = %v", err)
	}
}
