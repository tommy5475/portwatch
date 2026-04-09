package reporter

import (
	"strings"
	"testing"

	"portwatch/internal/state"
)

func TestFormatTextContainsProtocol(t *testing.T) {
	d := state.Diff{
		Timestamp: 1_700_000_000,
		Opened:    []state.Port{{Port: 8080, Protocol: "tcp"}},
		Closed:    []state.Port{{Port: 22, Protocol: "tcp"}},
	}
	out := formatText(d)
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in text output, got: %s", out)
	}
	if !strings.Contains(out, "tcp") {
		t.Errorf("expected protocol tcp in text output, got: %s", out)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("expected '+' marker for opened port, got: %s", out)
	}
	if !strings.Contains(out, "-") {
		t.Errorf("expected '-' marker for closed port, got: %s", out)
	}
}

func TestFormatCSVFields(t *testing.T) {
	d := state.Diff{
		Timestamp: 1_700_000_000,
		Opened:    []state.Port{{Port: 443, Protocol: "tcp"}},
		Closed:    []state.Port{},
	}
	out := formatCSV(d)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected header + at least one data row, got %d lines", len(lines))
	}
	if lines[0] != "timestamp,event,port,protocol" {
		t.Errorf("unexpected CSV header: %s", lines[0])
	}
	if !strings.Contains(lines[1], "opened") {
		t.Errorf("expected 'opened' in data row, got: %s", lines[1])
	}
	if !strings.Contains(lines[1], "443") {
		t.Errorf("expected port 443 in data row, got: %s", lines[1])
	}
}

func TestDefaultFormatIsText(t *testing.T) {
	f, err := ParseFormat("text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != FormatText {
		t.Errorf("expected FormatText, got %q", f)
	}
}

func TestParseFormatUnknown(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestFormatTextEmptyDiff(t *testing.T) {
	d := state.Diff{Timestamp: 1_700_000_000}
	out := formatText(d)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes' message for empty diff, got: %s", out)
	}
}
