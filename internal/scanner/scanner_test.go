package scanner

import (
	"net"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	timeout := 100 * time.Millisecond
	s := New(timeout)

	if s.timeout != timeout {
		t.Errorf("expected timeout %v, got %v", timeout, s.timeout)
	}
}

func TestScanTCPPort(t *testing.T) {
	s := New(100 * time.Millisecond)

	// Start a test server
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	// Test open port
	open, err := s.ScanTCPPort(port)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !open {
		t.Error("expected port to be open")
	}

	// Test closed port (use a very high port unlikely to be in use)
	open, err = s.ScanTCPPort(65534)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if open {
		t.Error("expected port to be closed")
	}
}

func TestScanTCPRange(t *testing.T) {
	s := New(100 * time.Millisecond)

	// Test invalid range
	_, err := s.ScanTCPRange(100, 50)
	if err == nil {
		t.Error("expected error for invalid range")
	}

	_, err = s.ScanTCPRange(0, 100)
	if err == nil {
		t.Error("expected error for port < 1")
	}

	_, err = s.ScanTCPRange(100, 70000)
	if err == nil {
		t.Error("expected error for port > 65535")
	}

	// Test valid range
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	openPorts, err := s.ScanTCPRange(port, port+10)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(openPorts) != 1 {
		t.Errorf("expected 1 open port, got %d", len(openPorts))
	}

	if len(openPorts) > 0 && openPorts[0].Number != port {
		t.Errorf("expected port %d, got %d", port, openPorts[0].Number)
	}
}
