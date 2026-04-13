package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestDrainAcquireRelease(t *testing.T) {
	d := newDrain()
	if !d.Acquire() {
		t.Fatal("expected Acquire to succeed on fresh drain")
	}
	if d.Inflight() != 1 {
		t.Fatalf("expected 1 inflight, got %d", d.Inflight())
	}
	d.Release()
	if d.Inflight() != 0 {
		t.Fatalf("expected 0 inflight after Release, got %d", d.Inflight())
	}
}

func TestDrainCloseBlocksNewAcquire(t *testing.T) {
	d := newDrain()
	go d.Close(100 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	if d.Acquire() {
		t.Fatal("expected Acquire to fail after Close")
	}
}

func TestDrainCloseReturnsTrueWhenEmpty(t *testing.T) {
	d := newDrain()
	ok := d.Close(50 * time.Millisecond)
	if !ok {
		t.Fatal("expected Close to return true when no inflight work")
	}
}

func TestDrainCloseWaitsForRelease(t *testing.T) {
	d := newDrain()
	d.Acquire()
	go func() {
		time.Sleep(20 * time.Millisecond)
		d.Release()
	}()
	ok := d.Close(200 * time.Millisecond)
	if !ok {
		t.Fatal("expected Close to return true after Release")
	}
}

func TestDrainCloseTimeoutExceeded(t *testing.T) {
	d := newDrain()
	d.Acquire() // never released
	ok := d.Close(20 * time.Millisecond)
	if ok {
		t.Fatal("expected Close to return false on timeout")
	}
}

func TestDrainConcurrentWorkers(t *testing.T) {
	d := newDrain()
	const workers = 50
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			if d.Acquire() {
				time.Sleep(5 * time.Millisecond)
				d.Release()
			}
		}()
	}
	wg.Wait()
	ok := d.Close(100 * time.Millisecond)
	if !ok {
		t.Fatal("expected drain to complete after all workers finished")
	}
	if d.Inflight() != 0 {
		t.Fatalf("expected 0 inflight, got %d", d.Inflight())
	}
}
