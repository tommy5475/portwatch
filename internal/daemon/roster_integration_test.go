package daemon

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRosterConcurrentJoinCheckinLeave(t *testing.T) {
	r := newRoster()
	const workers = 50
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			name := fmt.Sprintf("worker-%d", id)
			r.join(name)
			for j := 0; j < 5; j++ {
				r.checkin(name)
			}
			r.leave(name)
		}(i)
	}

	wg.Wait()
	if r.len() != 0 {
		t.Fatalf("expected empty roster after all workers left, got %d", r.len())
	}
}

func TestRosterMarkStaleConcurrentCheckin(t *testing.T) {
	r := newRoster()
	const workers = 20

	for i := 0; i < workers; i++ {
		r.join(fmt.Sprintf("w-%d", i))
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			r.markStale(time.Millisecond)
			time.Sleep(time.Microsecond)
		}
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				r.checkin(fmt.Sprintf("w-%d", id))
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	wg.Wait()
	// no panic or race is the success criterion
}
