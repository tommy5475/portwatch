package reporter_test

import (
	"strings"
	"testing"

	"portwatch/internal/reporter"
	"portwatch/internal/state"
)

func makeDiff() state.Diff {
	return state.Diff{
		Opened: []state.PortEntry{{Protocol: "tcp", Port: 8080}, {Protocol: "udp", Port: 53}},
		Closed: []state.PortEntry{{Protocol: "tcp", Port: 22}},
	}
}

func TestNew(t *testing.T) {
	r := reporter.New(nil, "")
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestReportText(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	diff := makeDiff()

	if err := r.Report(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "OPENED") {
		t.Error("expected OPENED in text output")
	}
	if !strings.Contains(out, "CLOSED") {
		t.Error("expected CLOSED in text output")
	}
	if !strings.Contains(out, "8080") {
		t.Error("expected port 8080 in output")
	}
	if !strings.Contains(out, "22") {
		t.Error("expected port 22 in output")
	}
}

func TestReportCSV(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatCSV)
	diff := makeDiff()

	if err := r.Report(diff); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 CSV lines, got %d", len(lines))
	}
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			t.Errorf("expected 4 CSV fields, got %d in line: %s", len(parts), line)
		}
	}
}

func TestReportEmptyDiff(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf, reporter.FormatText)
	empty := state.Diff{}

	if err := r.Report(empty); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "OPENED") || strings.Contains(out, "CLOSED") {
		t.Error("expected no change entries for empty diff")
	}
}
