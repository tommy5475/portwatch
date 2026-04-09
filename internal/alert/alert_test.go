package alert

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	alerter := New(nil)
	if alerter == nil {
		t.Fatal("Expected non-nil alerter")
	}
	if alerter.logger == nil {
		t.Fatal("Expected non-nil logger")
	}
}

func TestAlert(t *testing.T) {
	tests := []struct {
		name     string
		change   PortChange
		expected string
	}{
		{
			name: "port opened",
			change: PortChange{
				Port:      8080,
				OldState:  false,
				NewState:  true,
				Timestamp: time.Now(),
			},
			expected: "PORT OPENED: Port 8080 is now OPEN (was CLOSED)",
		},
		{
			name: "port closed",
			change: PortChange{
				Port:      3306,
				OldState:  true,
				NewState:  false,
				Timestamp: time.Now(),
			},
			expected: "PORT CLOSED: Port 3306 is now CLOSED (was OPEN)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tmpFile, _ := os.CreateTemp("", "alert_test")
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			alerter := New(tmpFile)
			alerter.logger.SetOutput(&buf)

			err := alerter.Alert(tt.change)
			if err != nil {
				t.Fatalf("Alert failed: %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestAlertBatch(t *testing.T) {
	var buf bytes.Buffer
	tmpFile, _ := os.CreateTemp("", "alert_test")
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	alerter := New(tmpFile)
	alerter.logger.SetOutput(&buf)

	changes := []PortChange{
		{Port: 80, OldState: false, NewState: true, Timestamp: time.Now()},
		{Port: 443, OldState: true, NewState: false, Timestamp: time.Now()},
	}

	err := alerter.AlertBatch(changes)
	if err != nil {
		t.Fatalf("AlertBatch failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Detected 2 port change(s)") {
		t.Errorf("Expected batch summary in output")
	}
	if !strings.Contains(output, "Port 80") {
		t.Errorf("Expected Port 80 in output")
	}
	if !strings.Contains(output, "Port 443") {
		t.Errorf("Expected Port 443 in output")
	}
}
