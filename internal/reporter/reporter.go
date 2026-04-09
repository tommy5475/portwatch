// Package reporter provides functionality for generating human-readable
// reports of port scan results and changes detected by the monitor.
package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"portwatch/internal/state"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatCSV  Format = "csv"
)

// Reporter writes port change reports to an output destination.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a new Reporter writing to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer, format Format) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Reporter{out: w, format: format}
}

// Report writes a formatted report of the given diff to the reporter's output.
func (r *Reporter) Report(diff state.Diff) error {
	switch r.format {
	case FormatCSV:
		return r.writeCSV(diff)
	default:
		return r.writeText(diff)
	}
}

func (r *Reporter) writeText(diff state.Diff) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] Port changes detected:\n", timestamp))
	for _, p := range diff.Opened {
		sb.WriteString(fmt.Sprintf("  + OPENED  %s/%d\n", strings.ToUpper(p.Protocol), p.Port))
	}
	for _, p := range diff.Closed {
		sb.WriteString(fmt.Sprintf("  - CLOSED  %s/%d\n", strings.ToUpper(p.Protocol), p.Port))
	}
	_, err := fmt.Fprint(r.out, sb.String())
	return err
}

func (r *Reporter) writeCSV(diff state.Diff) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	var sb strings.Builder
	for _, p := range diff.Opened {
		sb.WriteString(fmt.Sprintf("%s,opened,%s,%d\n", timestamp, p.Protocol, p.Port))
	}
	for _, p := range diff.Closed {
		sb.WriteString(fmt.Sprintf("%s,closed,%s,%d\n", timestamp, p.Protocol, p.Port))
	}
	_, err := fmt.Fprint(r.out, sb.String())
	return err
}
