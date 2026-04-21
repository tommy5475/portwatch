// Package daemon provides internal building blocks for the portwatch daemon.
//
// # Sampler
//
// sampler is a lightweight, goroutine-safe numeric statistics collector.
// It is designed for recording periodic measurements such as scan durations,
// port counts, or alert latencies.
//
// Usage:
//
//	s := newSampler()
//	s.record(42.5)
//	s.record(38.0)
//	stats := s.stats()
//	fmt.Println(stats.Mean) // 40.25
//
// All methods are safe for concurrent use.
package daemon
