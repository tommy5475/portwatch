package daemon

import (
	"fmt"
	"sync"
	"testing"
)

func TestRegistryConcurrentRegisterGet(t *testing.T) {
	r := newRegistry()
	const workers = 20
	const entries = 50

	var wg sync.WaitGroup

	// Concurrent writers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < entries; i++ {
				name := fmt.Sprintf("worker-%d-entry-%d", id, i)
				r.Register(name, id, "worker")
			}
		}(w)
	}

	// Concurrent readers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < entries; i++ {
				_ = r.Len()
				_ = r.Names()
				_ = r.FilterByTag("worker")
			}
		}()
	}

	wg.Wait()

	if r.Len() != workers*entries {
		t.Fatalf("expected %d entries, got %d", workers*entries, r.Len())
	}
}

func TestRegistryConcurrentUnregister(t *testing.T) {
	r := newRegistry()
	const n = 100
	for i := 0; i < n; i++ {
		r.Register(fmt.Sprintf("entry-%d", i), i)
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			r.Unregister(fmt.Sprintf("entry-%d", idx))
		}(i)
	}
	wg.Wait()

	if r.Len() != 0 {
		t.Fatalf("expected 0 entries after concurrent unregister, got %d", r.Len())
	}
}
