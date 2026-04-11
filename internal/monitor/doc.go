// Package monitor provides port monitoring functionality for detecting changes
// in port states over time.
//
// The monitor package continuously scans specified ports at regular intervals
// and reports any state changes (open -> closed or closed -> open) through
// a channel-based notification system.
//
// # Architecture
//
// The monitor uses a polling approach: at each interval tick, it checks all
// configured ports and compares the results against the previously recorded
// state. Any differences are emitted as [Change] events on the changes channel.
//
// # Example usage
//
//	package main
//
//	import (
//		"context"
//		"fmt"
//		"time"
//		"portwatch/internal/monitor"
//	)
//
//	func main() {
//		ports := []int{80, 443, 8080}
//		interval := 5 * time.Second
//
//		m := monitor.New(ports, interval)
//		ctx := context.Background()
//
//		// Listen for changes before starting to avoid missing events.
//		go func() {
//			for change := range m.Changes() {
//				fmt.Printf("Port %d changed: %v -> %v\n",
//					change.Port, change.WasOpen, change.NowOpen)
//			}
//		}()
//
//		// Start monitoring (blocks until ctx is cancelled or an error occurs).
//		if err := m.Start(ctx); err != nil {
//			fmt.Printf("Monitor error: %v\n", err)
//		}
//	}
//
// # Thread safety
//
// The monitor maintains thread-safe state tracking and provides methods to
// query the current state of all monitored ports. All exported methods are
// safe for concurrent use.
package monitor
