package filter

import "testing"

func TestBuilderDenyPort(t *testing.T) {
	f, err := NewBuilder().
		DenyPort("tcp", 22).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Allow("tcp", 22) {
		t.Error("expected port 22 to be denied")
	}
	if !f.Allow("tcp", 80) {
		t.Error("expected port 80 to be allowed")
	}
}

func TestBuilderAllowPort(t *testing.T) {
	f, err := NewBuilder().
		DenyRange("tcp", 1, 1024).
		AllowPort("tcp", 80).
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// First rule (deny 1-1024) matches port 443 before allow rule
	if f.Allow("tcp", 443) {
		t.Error("expected port 443 to be denied by range")
	}
}

func TestBuilderInvalidRange(t *testing.T) {
	_, err := NewBuilder().
		DenyRange("tcp", 500, 100).
		Build()
	if err == nil {
		t.Fatal("expected error for invalid range in builder")
	}
}

func TestBuilderChainStopsOnError(t *testing.T) {
	b := NewBuilder().
		DenyRange("tcp", 500, 100). // invalid — sets error
		AllowPort("tcp", 80)        // should be skipped
	if len(b.rules) != 0 {
		t.Error("rules should not accumulate after error")
	}
}

func TestBuilderEmpty(t *testing.T) {
	f, err := NewBuilder().Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow("tcp", 8080) {
		t.Error("empty filter should allow all")
	}
}
