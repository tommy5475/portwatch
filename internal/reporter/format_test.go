package reporter_test

import (
	"strings"
	"testing"

	"portwatch/internal/reporter"
	"portwatch/internal/state"
)

func TestFormatTextContainsProtocol(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	diff := state.Diff{
		Opened: []state.PortEntry{{Protocol: "udp", Port: 123}},
	}

	if err := r.Report(diff); err != nil {
		t.Fatalf("Report error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "UDP") {
		t.Errorf("expected protocol UDP in output, got: %s", out)
	}
	if !strings.Contains(out, "123") {
		t.Errorf("expected port 123 in output, got: %s", out)
	}
}

func TestFormatCSVFields(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatCSV)
	diff := state.Diff{
		Closed: []state.PortEntry{{Protocol: "tcp", Port: 443}},
	}

	if err := r.Report(diff); err != nil {
		t.Fatalf("Report error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	parts := strings.Split(line, ",")
	if len(parts) != 4 {
		t.Fatalf("expected 4 CSV fields, got %d: %s", len(parts), line)
	}
	if parts[1] != "closed" {
		t.Errorf("expected status 'closed', got %q", parts[1])
	}
	if parts[2] != "tcp" {
		t.Errorf("expected protocol 'tcp', got %q", parts[2])
	}
	if parts[3] != "443" {
		t.Errorf("expected port '443', got %q", parts[3])
	}
}

func TestDefaultFormatIsText(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, "")
	diff := state.Diff{
		Opened: []state.PortEntry{{Protocol: "tcp", Port: 9090}},
	}

	if err := r.Report(diff); err != nil {
		t.Fatalf("Report error: %v", err)
	}

	out := buf.String()
	// text format includes the word "changes"
	if !strings.Contains(out, "changes") {
		t.Errorf("expected text-format output, got: %s", out)
	}
}
