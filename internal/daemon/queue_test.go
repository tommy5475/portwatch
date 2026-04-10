package daemon

import (
	"testing"
)

func TestQueuePushPop(t *testing.T) {
	q := newQueue[int](4)
	q.Push(1)
	q.Push(2)

	v, ok := q.Pop()
	if !ok || v != 1 {
		t.Fatalf("expected 1, got %d ok=%v", v, ok)
	}
	v, ok = q.Pop()
	if !ok || v != 2 {
		t.Fatalf("expected 2, got %d ok=%v", v, ok)
	}
}

func TestQueuePopEmpty(t *testing.T) {
	q := newQueue[string](4)
	_, ok := q.Pop()
	if ok {
		t.Fatal("expected false on empty pop")
	}
}

func TestQueueLen(t *testing.T) {
	q := newQueue[int](8)
	for i := range 5 {
		q.Push(i)
	}
	if q.Len() != 5 {
		t.Fatalf("expected len 5, got %d", q.Len())
	}
}

func TestQueueEviction(t *testing.T) {
	q := newQueue[int](3)
	for i := range 5 {
		q.Push(i)
	}
	if q.Evicted() != 2 {
		t.Fatalf("expected 2 evictions, got %d", q.Evicted())
	}
	// oldest surviving item should be 2
	v, ok := q.Pop()
	if !ok || v != 2 {
		t.Fatalf("expected front=2, got %d ok=%v", v, ok)
	}
}

func TestQueueDrain(t *testing.T) {
	q := newQueue[int](8)
	q.Push(10)
	q.Push(20)
	q.Push(30)

	out := q.Drain()
	if len(out) != 3 {
		t.Fatalf("expected 3 items, got %d", len(out))
	}
	if q.Len() != 0 {
		t.Fatal("queue should be empty after drain")
	}
}

func TestQueueDefaultCapacity(t *testing.T) {
	q := newQueue[int](0)
	for i := range 64 {
		q.Push(i)
	}
	if q.Len() != 64 {
		t.Fatalf("expected 64 items, got %d", q.Len())
	}
	if q.Evicted() != 0 {
		t.Fatalf("expected 0 evictions, got %d", q.Evicted())
	}
}

func TestQueueFIFOOrder(t *testing.T) {
	q := newQueue[int](8)
	for i := range 5 {
		q.Push(i)
	}
	for i := range 5 {
		v, ok := q.Pop()
		if !ok || v != i {
			t.Fatalf("position %d: expected %d, got %d", i, i, v)
		}
	}
}
