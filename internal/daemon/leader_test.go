package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestLeaderInitiallyNotLeader(t *testing.T) {
	l := newLeader()
	if l.IsLeader() {
		t.Fatal("expected not leader initially")
	}
}

func TestLeaderAcquireMakesLeader(t *testing.T) {
	l := newLeader()
	l.Acquire()
	if !l.IsLeader() {
		t.Fatal("expected leader after Acquire")
	}
}

func TestLeaderReleaseClearsLeader(t *testing.T) {
	l := newLeader()
	l.Acquire()
	l.Release()
	if l.IsLeader() {
		t.Fatal("expected not leader after Release")
	}
}

func TestLeaderTermIncrementsOnEachAcquire(t *testing.T) {
	l := newLeader()
	if l.Term() != 0 {
		t.Fatalf("expected initial term 0, got %d", l.Term())
	}
	l.Acquire()
	if l.Term() != 1 {
		t.Fatalf("expected term 1, got %d", l.Term())
	}
	l.Acquire()
	if l.Term() != 2 {
		t.Fatalf("expected term 2, got %d", l.Term())
	}
}

func TestLeaderTermDoesNotResetOnRelease(t *testing.T) {
	l := newLeader()
	l.Acquire()
	l.Release()
	l.Acquire()
	if l.Term() != 2 {
		t.Fatalf("expected term 2 after re-acquire, got %d", l.Term())
	}
}

func TestLeaderAgeIsZeroWhenNotLeader(t *testing.T) {
	l := newLeader()
	if l.Age() != 0 {
		t.Fatal("expected zero age when not leader")
	}
}

func TestLeaderAgeIsPositiveWhenLeader(t *testing.T) {
	l := newLeader()
	l.Acquire()
	time.Sleep(2 * time.Millisecond)
	if l.Age() <= 0 {
		t.Fatal("expected positive age while leading")
	}
}

func TestLeaderLastRenewedUpdatesOnAcquire(t *testing.T) {
	l := newLeader()
	before := time.Now()
	l.Acquire()
	after := time.Now()
	rn := l.LastRenewed()
	if rn.Before(before) || rn.After(after) {
		t.Fatalf("LastRenewed %v not between %v and %v", rn, before, after)
	}
}

func TestLeaderConcurrentAccess(t *testing.T) {
	l := newLeader()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				l.Acquire()
			} else {
				l.Release()
			}
			_ = l.IsLeader()
			_ = l.Term()
			_ = l.Age()
		}(i)
	}
	wg.Wait()
}
