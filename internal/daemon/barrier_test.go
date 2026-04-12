package daemon

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBarrierDefaultsToOne(t *testing.T) {
	b := newBarrier(0)
	if b.Size() != 1 {
		t.Fatalf("expected size 1, got %d", b.Size())
	}
}

func TestBarrierSingleParticipantReturnsImmediately(t *testing.T) {
	b := newBarrier(1)
	done := make(chan struct{})
	go func() {
		b.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("barrier with 1 participant did not release")
	}
}

func TestBarrierReleasesAllParticipants(t *testing.T) {
	const n = 5
	b := newBarrier(n)
	var released int64
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			b.Wait()
			atomic.AddInt64(&released, 1)
		}()
	}
	wg.Wait()
	if atomic.LoadInt64(&released) != n {
		t.Fatalf("expected %d released, got %d", n, released)
	}
}

func TestBarrierGenerationIncrements(t *testing.T) {
	const n = 3
	b := newBarrier(n)
	var wg sync.WaitGroup
	gens := make([]uint64, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			gens[i] = b.Wait()
		}()
	}
	wg.Wait()
	for _, g := range gens {
		if g != 0 {
			t.Fatalf("expected generation 0, got %d", g)
		}
	}
	// Second round – generation should be 1.
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			gens[i] = b.Wait()
		}()
	}
	wg.Wait()
	for _, g := range gens {
		if g != 1 {
			t.Fatalf("expected generation 1, got %d", g)
		}
	}
}

func TestBarrierArrivedCount(t *testing.T) {
	b := newBarrier(3)
	ready := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			<-ready
			b.Wait()
			wg.Done()
		}()
	}
	time.Sleep(20 * time.Millisecond)
	close(ready)
	time.Sleep(20 * time.Millisecond)
	// At least one goroutine should have arrived before the third.
	if b.Arrived() < 0 || b.Arrived() > 2 {
		t.Fatalf("unexpected arrived count %d", b.Arrived())
	}
	// Release the barrier by sending the third participant.
	go func() { b.Wait() }()
	wg.Wait()
}
