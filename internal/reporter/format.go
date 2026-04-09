package reporter

import (
	"fmt"
	"strings"
	"time"

	"portwatch/internal/state"
)

// Format represents an output format for reports.
type Format string

const (
	// FormatText renders diffs as human-readable plain text.
	FormatText Format = "text"
	// FormatCSV renders diffs as comma-separated values.
	FormatCSV Format = "csv"
)

// ParseFormat converts a string to a Format, returning an error if unknown.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case string(FormatText):
		return FormatText, nil
	case string(FormatCSV):
		return FormatCSV, nil
	default:
		return "", fmt.Errorf("unknown format %q: must be \"text\" or \"csv\"", s)
	}
}

// formatText renders a state.Diff as a human-readable string.
func formatText(d state.Diff) string {
	if len(d.Opened)+len(d.Closed) == 0 {
		return "no changes detected\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("port changes at %s\n", time.Unix(d.Timestamp, 0).UTC().Format(time.RFC3339)))
	for _, p := range d.Opened {
		sb.WriteString(fmt.Sprintf("  + %-6d %s\n", p.Port, p.Protocol))
	}
	for _, p := range d.Closed {
		sb.WriteString(fmt.Sprintf("  - %-6d %s\n", p.Port, p.Protocol))
	}
	return sb.String()
}

// formatCSV renders a state.Diff as CSV rows (header + data).
func formatCSV(d state.Diff) string {
	var sb strings.Builder
	sb.WriteString("timestamp,event,port,protocol\n")
	ts := time.Unix(d.Timestamp, 0).UTC().Format(time.RFC3339)
	for _, p := range d.Opened {
		sb.WriteString(fmt.Sprintf("%s,opened,%d,%s\n", ts, p.Port, p.Protocol))
	}
	for _, p := range d.Closed {
		sb.WriteString(fmt.Sprintf("%s,closed,%d,%s\n", ts, p.Port, p.Protocol))
	}
	return sb.String()
}
