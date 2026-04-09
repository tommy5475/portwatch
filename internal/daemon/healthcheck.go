package daemon

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"
)

// HealthStatus represents the current health of the daemon.
type HealthStatus struct {
	Status    string    `json:"status"`
	Uptime    string    `json:"uptime"`
	Scans     uint64    `json:"scans_completed"`
	StartedAt time.Time `json:"started_at"`
}

// healthServer exposes a minimal HTTP health endpoint.
type healthServer struct {
	addr      string
	scans     atomic.Uint64
	startedAt time.Time
	server    *http.Server
}

func newHealthServer(addr string) *healthServer {
	hs := &healthServer{
		addr:      addr,
		startedAt: time.Now(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", hs.handleHealth)

	hs.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return hs
}

func (hs *healthServer) start() error {
	go func() { _ = hs.server.ListenAndServe() }()
	return nil
}

func (hs *healthServer) stop() {
	_ = hs.server.Close()
}

func (hs *healthServer) recordScan() {
	hs.scans.Add(1)
}

func (hs *healthServer) handleHealth(w http.ResponseWriter, _ *http.Request) {
	status := HealthStatus{
		Status:    "ok",
		Uptime:    time.Since(hs.startedAt).Round(time.Second).String(),
		Scans:     hs.scans.Load(),
		StartedAt: hs.startedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
