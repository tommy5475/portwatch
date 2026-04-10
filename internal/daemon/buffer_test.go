package daemon

import (
	"sync"
	"testing"
)

func TestBufferEmptySnapshot(t *testing.T) {
	b := newBuffer[int](4)
	if got := b.Snapshot(); got != nil {
		t.Fatalf("expected nil snapshot on empty buffer, got %v", got)
	}
}

func TestBufferLenAfterPush(t *testing.T) {
	b := newBuffer[int](4)
	b.Push(1)
	b.Push(2)
	if b.Len() != 2 {
		t.Fatalf("expected len 2, got %d", b.Len())
	}
}

func TestBufferSnapshotOrder(t *testing.T) {
	b := newBuffer[int](4)
	for _, v := range []int{10, 20, 30} {
		b.Push(v)
	}
	snap := b.Snapshot()
	expected := []int{10, 20, 30}
	for i, v := range expected {
		if snap[i] != v {
			t.Fatalf("index %d: expected %d, got %d", i, v, snap[i])
		}
	}
}

func TestBufferOverwrite(t *testing.T) {
	b := newBuffer[int](3)
	for _, v := range []int{1, 2, 3, 4} {
		b.Push(v)
	}
	snap := b.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 items, got %d", len(snap))
	}
	if snap[0] != 2 || snap[1] != 3 || snap[2] != 4 {
		t.Fatalf("unexpected snapshot after overwrite: %v", snap)
	}
}

func TestBufferReset(t *testing.T) {
	b := newBuffer[int](4)
	b.Push(1)
	b.Push(2)
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected len 0 after reset, got %d", b.Len())
	}
	if b.Snapshot() != nil {
		t.Fatal("expected nil snapshot after reset")
	}
}

func TestBufferDefaultCapacity(t *testing.T) {
	b := newBuffer[string](0)
	for i := 0; i < 20; i++ {
		b.Push("x")
	}
	if b.Len() != 16 {
		t.Fatalf("expected default cap 16, got %d", b.Len())
	}
}

func TestBufferConcurrentAccess(t *testing.T) {
	b := newBuffer[int](64)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			b.Push(v)
			_ = b.Snapshot()
		}(i)
	}
	wg.Wait()
	if b.Len() == 0 {
		t.Fatal("expected items after concurrent pushes")
	}
}
