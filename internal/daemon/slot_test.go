package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestSlotAcquireSucceeds(t *testing.T) {
	s := newSlot("test", time.Second)
	if !s.Acquire() {
		t.Fatal("expected first acquire to succeed")
	}
}

func TestSlotAcquireBlocksWhenHeld(t *testing.T) {
	s := newSlot("test", time.Second)
	s.Acquire()
	if s.Acquire() {
		t.Fatal("expected second acquire to fail while held")
	}
}

func TestSlotReleaseAllowsReacquire(t *testing.T) {
	s := newSlot("test", time.Second)
	s.Acquire()
	s.Release()
	if !s.Acquire() {
		t.Fatal("expected acquire to succeed after release")
	}
}

func TestSlotHeldAfterAcquire(t *testing.T) {
	s := newSlot("test", time.Second)
	s.Acquire()
	if !s.Held() {
		t.Fatal("expected slot to be held after acquire")
	}
}

func TestSlotNotHeldAfterRelease(t *testing.T) {
	s := newSlot("test", time.Second)
	s.Acquire()
	s.Release()
	if s.Held() {
		t.Fatal("expected slot not held after release")
	}
}

func TestSlotTTLExpiry(t *testing.T) {
	s := newSlot("test", 10*time.Millisecond)
	s.Acquire()
	time.Sleep(20 * time.Millisecond)
	if !s.Acquire() {
		t.Fatal("expected acquire to succeed after TTL expiry")
	}
}

func TestSlotCountIncrements(t *testing.T) {
	s := newSlot("test", time.Second)
	for i := 0; i < 3; i++ {
		s.Acquire()
		s.Release()
	}
	if s.Count() != 3 {
		t.Fatalf("expected count 3, got %d", s.Count())
	}
}

func TestSlotDefaultsOnInvalidArgs(t *testing.T) {
	s := newSlot("", 0)
	if s.Name() != "default" {
		t.Fatalf("expected default name, got %q", s.Name())
	}
	if !s.Acquire() {
		t.Fatal("expected acquire to succeed")
	}
}

func TestSlotConcurrentAcquire(t *testing.T) {
	s := newSlot("test", time.Second)
	var wg sync.WaitGroup
	wins := make(chan struct{}, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if s.Acquire() {
				wins <- struct{}{}
			}
		}()
	}
	wg.Wait()
	close(wins)
	count := 0
	for range wins {
		count++
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 winner, got %d", count)
	}
}
