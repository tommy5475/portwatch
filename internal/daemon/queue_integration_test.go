package daemon

import (
	"sync"
	"testing"
)

func TestQueueConcurrentPushPop(t *testing.T) {
	const workers = 8
	const perWorker = 200

	q := newQueue[int](workers * perWorker)
	var wg sync.WaitGroup

	// concurrent producers
	for w := range workers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := range perWorker {
				q.Push(id*perWorker + i)
			}
		}(w)
	}
	wg.Wait()

	if q.Len() != workers*perWorker {
		t.Fatalf("expected %d items, got %d", workers*perWorker, q.Len())
	}
}

func TestQueueConcurrentEviction(t *testing.T) {
	const cap = 50
	const total = 500

	q := newQueue[int](cap)
	var wg sync.WaitGroup

	for i := range total {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			q.Push(v)
		}(i)
	}
	wg.Wait()

	if q.Len() > cap {
		t.Fatalf("queue exceeded capacity: len=%d cap=%d", q.Len(), cap)
	}
	if q.Evicted() == 0 {
		t.Fatal("expected evictions under overflow, got none")
	}
}
