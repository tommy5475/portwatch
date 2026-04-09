package filter

import (
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New([]Rule{
		{Protocol: "tcp", PortMin: 80, PortMax: 80, Allow: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewInvalidRange(t *testing.T) {
	_, err := New([]Rule{
		{Protocol: "tcp", PortMin: 100, PortMax: 80, Allow: true},
	})
	if err == nil {
		t.Fatal("expected error for invalid range")
	}
}

func TestNewUnknownProtocol(t *testing.T) {
	_, err := New([]Rule{
		{Protocol: "sctp", PortMin: 80, PortMax: 80, Allow: true},
	})
	if err == nil {
		t.Fatal("expected error for unknown protocol")
	}
}

func TestAllowDefaultPermit(t *testing.T) {
	f, _ := New(nil)
	if !f.Allow("tcp", 8080) {
		t.Error("expected default allow")
	}
}

func TestAllowExplicitDeny(t *testing.T) {
	f, _ := New([]Rule{
		{Protocol: "tcp", PortMin: 22, PortMax: 22, Allow: false},
	})
	if f.Allow("tcp", 22) {
		t.Error("expected port 22 to be denied")
	}
	if !f.Allow("tcp", 80) {
		t.Error("expected port 80 to be allowed")
	}
}

func TestAllowRangeRule(t *testing.T) {
	f, _ := New([]Rule{
		{Protocol: "*", PortMin: 1024, PortMax: 65535, Allow: false},
		{Protocol: "tcp", PortMin: 8080, PortMax: 8080, Allow: true},
	})
	// First matching rule wins
	if f.Allow("tcp", 9000) {
		t.Error("expected port 9000 to be denied by range rule")
	}
	if !f.Allow("tcp", 80) {
		t.Error("expected port 80 to be allowed (no rule matches)")
	}
}

func TestAllowProtocolIsolation(t *testing.T) {
	f, _ := New([]Rule{
		{Protocol: "udp", PortMin: 53, PortMax: 53, Allow: false},
	})
	if !f.Allow("tcp", 53) {
		t.Error("tcp/53 should not be affected by udp rule")
	}
	if f.Allow("udp", 53) {
		t.Error("udp/53 should be denied")
	}
}

func TestRulesCopy(t *testing.T) {
	input := []Rule{{Protocol: "tcp", PortMin: 80, PortMax: 80, Allow: true}}
	f, _ := New(input)
	copy := f.Rules()
	copy[0].PortMin = 999
	if f.Rules()[0].PortMin != 80 {
		t.Error("Rules() should return a copy, not a reference")
	}
}
