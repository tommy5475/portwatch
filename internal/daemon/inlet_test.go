package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestInletAcceptsWithinCapacity(t *testing.T) {
	in := newInlet[int](4)
	for i := 0; i < 4; i++ {
		if !in.Send(i) {
			t.Fatalf("expected Send to succeed for item %d", i)
		}
	}
	if in.Accepted() != 4 {
		t.Fatalf("expected accepted=4, got %d", in.Accepted())
	}
}

func TestInletDropsWhenFull(t *testing.T) {
	in := newInlet[int](2)
	in.Send(1)
	in.Send(2)
	if in.Send(3) {
		t.Fatal("expected Send to fail when inlet is full")
	}
	if in.Dropped() != 1 {
		t.Fatalf("expected dropped=1, got %d", in.Dropped())
	}
}

func TestInletLen(t *testing.T) {
	in := newInlet[string](8)
	in.Send("a")
	in.Send("b")
	if in.Len() != 2 {
		t.Fatalf("expected len=2, got %d", in.Len())
	}
}

func TestInletOut(t *testing.T) {
	in := newInlet[int](4)
	in.Send(42)
	select {
	case v := <-in.Out():
		if v != 42 {
			t.Fatalf("expected 42, got %d", v)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out reading from inlet")
	}
}

func TestInletDefaultCapacity(t *testing.T) {
	in := newInlet[int](0) // invalid — should default to 16
	for i := 0; i < 16; i++ {
		if !in.Send(i) {
			t.Fatalf("expected Send to succeed for item %d with default capacity", i)
		}
	}
	if in.Dropped() != 0 {
		t.Fatalf("expected no drops with default capacity, got %d", in.Dropped())
	}
}

func TestInletAgeIsPositive(t *testing.T) {
	in := newInlet[int](4)
	time.Sleep(2 * time.Millisecond)
	if in.Age() <= 0 {
		t.Fatal("expected positive age")
	}
}

func TestInletConcurrentSend(t *testing.T) {
	in := newInlet[int](64)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			in.Send(v)
		}(i)
	}
	wg.Wait()
	total := in.Accepted() + in.Dropped()
	if total != 100 {
		t.Fatalf("expected accepted+dropped=100, got %d", total)
	}
}
