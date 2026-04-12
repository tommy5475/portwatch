package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestCooldownConcurrentMarkAndReady(t *testing.T) {
	cd := newCooldown(50 * time.Millisecond)
	const workers = 20
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			cd.Ready("shared")
			cd.Mark("shared")
			cd.Remaining("shared")
		}()
	}
	wg.Wait()
	// No race — verified by -race flag.
}

func TestCooldownRealTimeExpiry(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real-time test in short mode")
	}
	cd := newCooldown(80 * time.Millisecond)
	cd.Mark("evt")
	if cd.Ready("evt") {
		t.Fatal("should not be ready immediately after mark")
	}
	time.Sleep(100 * time.Millisecond)
	if !cd.Ready("evt") {
		t.Fatal("should be ready after cooldown period")
	}
}
