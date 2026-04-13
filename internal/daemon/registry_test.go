package daemon

import (
	"sort"
	"testing"
)

func TestRegistryRegisterAndGet(t *testing.T) {
	r := newRegistry()
	r.Register("foo", 42)
	v, err := r.Get("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.(int) != 42 {
		t.Fatalf("expected 42, got %v", v)
	}
}

func TestRegistryGetMissing(t *testing.T) {
	r := newRegistry()
	_, err := r.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRegistryHas(t *testing.T) {
	r := newRegistry()
	if r.Has("x") {
		t.Fatal("expected false before register")
	}
	r.Register("x", "val")
	if !r.Has("x") {
		t.Fatal("expected true after register")
	}
}

func TestRegistryUnregister(t *testing.T) {
	r := newRegistry()
	r.Register("a", 1)
	if !r.Unregister("a") {
		t.Fatal("expected true for existing key")
	}
	if r.Has("a") {
		t.Fatal("key should be gone after unregister")
	}
	if r.Unregister("a") {
		t.Fatal("expected false for already-removed key")
	}
}

func TestRegistryLen(t *testing.T) {
	r := newRegistry()
	if r.Len() != 0 {
		t.Fatal("expected 0")
	}
	r.Register("a", 1)
	r.Register("b", 2)
	if r.Len() != 2 {
		t.Fatalf("expected 2, got %d", r.Len())
	}
}

func TestRegistryNames(t *testing.T) {
	r := newRegistry()
	r.Register("alpha", 1)
	r.Register("beta", 2)
	names := r.Names()
	sort.Strings(names)
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Fatalf("unexpected names: %v", names)
	}
}

func TestRegistryFilterByTag(t *testing.T) {
	r := newRegistry()
	r.Register("scanner", nil, "core", "active")
	r.Register("alerter", nil, "core")
	r.Register("debug", nil, "optional")

	core := r.FilterByTag("core")
	sort.Strings(core)
	if len(core) != 2 || core[0] != "alerter" || core[1] != "scanner" {
		t.Fatalf("unexpected core entries: %v", core)
	}

	opt := r.FilterByTag("optional")
	if len(opt) != 1 || opt[0] != "debug" {
		t.Fatalf("unexpected optional entries: %v", opt)
	}

	none := r.FilterByTag("unknown")
	if len(none) != 0 {
		t.Fatalf("expected empty, got %v", none)
	}
}

func TestRegistryOverwriteEntry(t *testing.T) {
	r := newRegistry()
	r.Register("key", "first")
	r.Register("key", "second")
	v, _ := r.Get("key")
	if v.(string) != "second" {
		t.Fatalf("expected second, got %v", v)
	}
	if r.Len() != 1 {
		t.Fatal("overwrite should not increase len")
	}
}
