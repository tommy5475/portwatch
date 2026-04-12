package daemon

import (
	"sync"
	"testing"
)

func TestRingbufEmpty(t *testing.T) {
	r := newRingbuf[int](4)
	if r.Len() != 0 {
		t.Fatalf("expected 0, got %d", r.Len())
	}
	if s := r.Snapshot(); s != nil {
		t.Fatalf("expected nil snapshot, got %v", s)
	}
}

func TestRingbufPushAndSnapshot(t *testing.T) {
	r := newRingbuf[int](4)
	for i := 1; i <= 3; i++ {
		r.Push(i)
	}
	s := r.Snapshot()
	if len(s) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(s))
	}
	for i, v := range s {
		if v != i+1 {
			t.Errorf("index %d: expected %d, got %d", i, i+1, v)
		}
	}
}

func TestRingbufOverwrite(t *testing.T) {
	r := newRingbuf[int](3)
	r.Push(1)
	r.Push(2)
	r.Push(3)
	r.Push(4) // overwrites 1
	s := r.Snapshot()
	if len(s) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(s))
	}
	expected := []int{2, 3, 4}
	for i, v := range s {
		if v != expected[i] {
			t.Errorf("index %d: expected %d, got %d", i, expected[i], v)
		}
	}
}

func TestRingbufReset(t *testing.T) {
	r := newRingbuf[int](4)
	r.Push(1)
	r.Push(2)
	r.Reset()
	if r.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", r.Len())
	}
	if s := r.Snapshot(); s != nil {
		t.Fatalf("expected nil after reset, got %v", s)
	}
}

func TestRingbufDefaultCap(t *testing.T) {
	r := newRingbuf[string](0)
	for i := 0; i < defaultRingbufCap+10; i++ {
		r.Push("x")
	}
	if r.Len() != defaultRingbufCap {
		t.Fatalf("expected %d, got %d", defaultRingbufCap, r.Len())
	}
}

func TestRingbufConcurrentAccess(t *testing.T) {
	r := newRingbuf[int](16)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			r.Push(v)
			_ = r.Snapshot()
			_ = r.Len()
		}(i)
	}
	wg.Wait()
	if r.Len() > 16 {
		t.Fatalf("len exceeded capacity: %d", r.Len())
	}
}
