package daemon

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBudgetConcurrentSpend(t *testing.T) {
	const cap = 50
	b := newBudget(cap, time.Minute) // long rate so no replenish during test

	var allowed atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if b.Spend(1) {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := allowed.Load(); got != cap {
		t.Fatalf("expected exactly %d allowed spends, got %d", cap, got)
	}
	if rem := b.Remaining(); rem != 0 {
		t.Fatalf("expected 0 remaining after exhaustion, got %d", rem)
	}
}

func TestBudgetReplenishUnderConcurrentLoad(t *testing.T) {
	b := newBudget(5, 10*time.Millisecond)
	b.Spend(5)

	var wg sync.WaitGroup
	var successes atomic.Int64

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(15 * time.Millisecond)
			if b.Spend(1) {
				successes.Add(1)
			}
		}()
	}
	wg.Wait()

	if successes.Load() == 0 {
		t.Fatal("expected at least one successful spend after replenish")
	}
}
