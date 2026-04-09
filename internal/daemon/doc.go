// Package daemon provides the top-level orchestration layer for portwatch.
//
// A Daemon is constructed from a [config.Config] and wires together every
// subsystem:
//
//   - [scanner.Scanner]  — low-level TCP/UDP port probing
//   - [monitor.Monitor]  — periodic scanning and change detection
//   - [filter.Filter]    — port/protocol allow-deny rules
//   - [reporter.Reporter] — human-readable or CSV diff formatting
//   - [notifier.Notifier] — stdout and optional webhook delivery
//
// Typical usage:
//
//	cfg, _ := config.LoadFromFile("/etc/portwatch/config.yaml")
//	d, _ := daemon.New(cfg)
//	if err := d.Run(ctx); err != nil && err != context.Canceled {
//		log.Fatal(err)
//	}
//
// The daemon runs until the supplied [context.Context] is cancelled, making
// it straightforward to integrate with OS signal handling.
package daemon
