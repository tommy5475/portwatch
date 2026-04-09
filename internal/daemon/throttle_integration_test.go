package daemon

import (
	"sync"
	"testing"
	"time"
)

// TestThrottleConcurrentAccess verifies that Allow and Reset are safe to call
// from multiple goroutines simultaneously.
func TestThrottleConcurrentAccess(t *testing.T) {
	th := newThrottle(10 * time.Millisecond)
	const workers = 20
	const iterations = 50

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			key := "key"
			for j := 0; j < iterations; j++ {
				th.Allow(key)
				if j%10 == 0 {
					th.Reset(key)
				}
				_ = th.Len()
			}
		}(i)
	}

	wg.Wait()
	// No race condition should be detected by the race detector.
}

// TestThrottleMultiKeyStorm simulates an alert storm across many distinct keys
// and confirms each key is individually gated.
func TestThrottleMultiKeyStorm(t *testing.T) {
	th := newThrottle(1 * time.Hour)
	keys := []string{"tcp:80", "tcp:443", "udp:53", "tcp:22", "tcp:3306"}

	allowed := 0
	for _, k := range keys {
		if th.Allow(k) {
			allowed++
		}
	}
	if allowed != len(keys) {
		t.Fatalf("expected %d allowed, got %d", len(keys), allowed)
	}

	// Second pass — all should be blocked.
	for _, k := range keys {
		if th.Allow(k) {
			t.Fatalf("key %q should be throttled", k)
		}
	}
}
