package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type healthStatus struct {
	Status    string `json:"status"`
	Uptime    string `json:"uptime"`
	Scans     int64  `json:"scans_total"`
	LastScan  string `json:"last_scan,omitempty"`
	Degraded  bool   `json:"degraded"`
}

type healthServer struct {
	start   time.Time
	scans   atomic.Int64
	lastTS  atomic.Value // stores time.Time
	watcher *watcher
	server  *http.Server
}

func newHealthServer(addr string, w *watcher) *healthServer {
	hs := &healthServer{
		start:   time.Now(),
		watcher: w,
	}
	hs.lastTS.Store(time.Time{})

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", hs.handleHealth)
	hs.server = &http.Server{Addr: addr, Handler: mux}
	return hs
}

func (hs *healthServer) recordScan() {
	hs.scans.Add(1)
	hs.lastTS.Store(time.Now())
}

func (hs *healthServer) handleHealth(w http.ResponseWriter, _ *http.Request) {
	degraded := hs.watcher != nil && hs.watcher.isDegraded()
	status := "ok"
	code := http.StatusOK
	if degraded {
		status = "degraded"
		code = http.StatusServiceUnavailable
	}

	lastScan := ""
	if ts, ok := hs.lastTS.Load().(time.Time); ok && !ts.IsZero() {
		lastScan = ts.UTC().Format(time.RFC3339)
	}

	body := healthStatus{
		Status:   status,
		Uptime:   fmt.Sprintf("%.0fs", time.Since(hs.start).Seconds()),
		Scans:    hs.scans.Load(),
		LastScan: lastScan,
		Degraded: degraded,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func (hs *healthServer) start_(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() { errCh <- hs.server.ListenAndServe() }()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		_ = hs.server.Shutdown(context.Background())
		return nil
	}
}
