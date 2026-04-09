// Package reporter provides human-readable and machine-parseable reporting
// of port state changes detected by portwatch.
//
// # Overview
//
// The reporter package formats [state.Diff] values into output suitable for
// display in a terminal or ingestion by external tooling. Two formats are
// supported:
//
//   - text  – a timestamped, human-readable summary (default)
//   - csv   – a comma-separated values stream suitable for log aggregation
//
// # Usage
//
//	r := reporter.New(os.Stdout, reporter.FormatText)
//	if err := r.Report(diff); err != nil {
//	    log.Fatal(err)
//	}
//
// Passing nil as the writer defaults to os.Stdout. Passing an empty format
// string defaults to [FormatText].
package reporter
