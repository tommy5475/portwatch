// Package monitor provides port monitoring functionality for detecting changes
// in port states over time.
//
// The monitor package continuously scans specified ports at regular intervals
// and reports any state changes (open -> closed or closed -> open) through
// a channel-based notification system.
//
// Example usage:
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
//		// Listen for changes
//		go func() {
//			for change := range m.Changes() {
//				fmt.Printf("Port %d changed: %v -> %v\n",
//					change.Port, change.WasOpen, change.NowOpen)
//			}
//		}()
//
//		// Start monitoring
//		if err := m.Start(ctx); err != nil {
//			fmt.Printf("Monitor error: %v\n", err)
//		}
//	}
//
// The monitor maintains thread-safe state tracking and provides methods to
// query the current state of all monitored ports.
package monitor
