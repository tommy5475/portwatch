package daemon

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

// freePort returns an available TCP port on localhost.
func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not find free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestHealthServerStartStop(t *testing.T) {
	port := freePort(t)
	addr := fmt.Sprintf("127.0.0.1:%d", port)

	hs := newHealthServer(addr)
	if err := hs.start(); err != nil {
		t.Fatalf("start failed: %v", err)
	}
	t.Cleanup(hs.stop)

	// Allow the server a moment to bind.
	time.Sleep(30 * time.Millisecond)

	hs.recordScan()

	url := fmt.Sprintf("http://%s/healthz", addr)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var status HealthStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if status.Scans != 1 {
		t.Errorf("expected 1 scan, got %d", status.Scans)
	}
}
