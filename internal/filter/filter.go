// Package filter provides port filtering logic for portwatch.
// It allows users to define rules that include or exclude specific
// ports or port ranges from monitoring.
package filter

import "fmt"

// Rule represents a single filter rule.
type Rule struct {
	Protocol string // "tcp", "udp", or "*" for both
	PortMin  uint16
	PortMax  uint16
	Allow    bool // true = include, false = exclude
}

// Filter holds an ordered list of rules applied to port checks.
type Filter struct {
	rules []Rule
}

// New returns a Filter with the given rules.
func New(rules []Rule) (*Filter, error) {
	for _, r := range rules {
		if r.PortMin > r.PortMax {
			return nil, fmt.Errorf("filter: invalid range %d-%d", r.PortMin, r.PortMax)
		}
		if r.Protocol != "tcp" && r.Protocol != "udp" && r.Protocol != "*" {
			return nil, fmt.Errorf("filter: unknown protocol %q", r.Protocol)
		}
	}
	return &Filter{rules: rules}, nil
}

// Allow reports whether the given protocol/port combination passes
// the filter. If no rule matches, the port is allowed by default.
func (f *Filter) Allow(protocol string, port uint16) bool {
	for _, r := range f.rules {
		if !matchesProtocol(r.Protocol, protocol) {
			continue
		}
		if port >= r.PortMin && port <= r.PortMax {
			return r.Allow
		}
	}
	return true // default allow
}

// Rules returns a copy of the current rule list.
func (f *Filter) Rules() []Rule {
	out := make([]Rule, len(f.rules))
	copy(out, f.rules)
	return out
}

func matchesProtocol(ruleProto, queryProto string) bool {
	return ruleProto == "*" || ruleProto == queryProto
}
