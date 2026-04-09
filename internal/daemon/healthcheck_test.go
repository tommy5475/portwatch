package daemon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthStatusFields(t *testing.T) {
	hs := newHealthServer("127.0.0.1:0")
	hs.recordScan()
	hs.recordScan()

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	hs.handleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var status HealthStatus
	if err := json.NewDecoder(rr.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if status.Status != "ok" {
		t.Errorf("expected status ok, got %q", status.Status)
	}
	if status.Scans != 2 {
		t.Errorf("expected 2 scans, got %d", status.Scans)
	}
	if status.StartedAt.IsZero() {
		t.Error("expected non-zero StartedAt")
	}
}

func TestHealthUptimeNonEmpty(t *testing.T) {
	hs := newHealthServer("127.0.0.1:0")
	// Simulate a small delay so uptime is non-trivial.
	time.Sleep(10 * time.Millisecond)

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	hs.handleHealth(rr, req)

	var status HealthStatus
	_ = json.NewDecoder(rr.Body).Decode(&status)

	if status.Uptime == "" {
		t.Error("expected non-empty uptime")
	}
}

func TestHealthContentType(t *testing.T) {
	hs := newHealthServer("127.0.0.1:0")
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	hs.handleHealth(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}
}

func TestRecordScanAtomic(t *testing.T) {
	hs := newHealthServer("127.0.0.1:0")
	const n = 50
	done := make(chan struct{})
	for i := 0; i < n; i++ {
		go func() { hs.recordScan(); done <- struct{}{} }()
	}
	for i := 0; i < n; i++ {
		<-done
	}
	if got := hs.scans.Load(); got != n {
		t.Errorf("expected %d scans, got %d", n, got)
	}
}
