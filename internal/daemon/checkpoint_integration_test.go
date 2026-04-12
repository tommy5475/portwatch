package daemon

import (
	"sync"
	"testing"
)

func TestCheckpointConcurrentMark(t *testing.T) {
	cp := newCheckpoint()
	const goroutines = 50
	const marksEach = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < marksEach; j++ {
				cp.Mark("event")
			}
		}()
	}
	wg.Wait()

	want := int64(goroutines * marksEach)
	if got := cp.Count("event"); got != want {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

func TestCheckpointConcurrentMultiKey(t *testing.T) {
	cp := newCheckpoint()
	keys := []string{"scan.start", "scan.complete", "alert.sent"}
	const goroutines = 20
	const marksEach = 50

	var wg sync.WaitGroup
	wg.Add(goroutines * len(keys))
	for _, k := range keys {
		for i := 0; i < goroutines; i++ {
			go func(key string) {
				defer wg.Done()
				for j := 0; j < marksEach; j++ {
					cp.Mark(key)
				}
			}(k)
		}
	}
	wg.Wait()

	want := int64(goroutines * marksEach)
	for _, k := range keys {
		if got := cp.Count(k); got != want {
			t.Errorf("key %q: expected %d, got %d", k, want, got)
		}
	}
}
