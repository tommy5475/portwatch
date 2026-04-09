// Package filter implements port-level allow/deny filtering for portwatch.
//
// # Overview
//
// A [Filter] holds an ordered list of [Rule] values. When portwatch scans
// a port, it consults the filter before deciding whether to track or alert
// on that port. The first matching rule wins; if no rule matches, the port
// is allowed by default.
//
// # Rules
//
// Each Rule specifies:
//   - Protocol: "tcp", "udp", or "*" (matches both)
//   - PortMin / PortMax: inclusive port range (single port when equal)
//   - Allow: whether matching ports are permitted (true) or suppressed (false)
//
// # Builder
//
// [Builder] provides a fluent API for constructing filters without
// manually populating Rule slices:
//
//	f, err := filter.NewBuilder().
//	    DenyPort("tcp", 22).
//	    AllowRange("tcp", 8000, 9000).
//	    Build()
//
// # Integration
//
// Filters are wired into the monitor via [config.Config]. Ports rejected
// by the filter are silently skipped during each scan cycle and will never
// appear in state snapshots or change notifications.
package filter
