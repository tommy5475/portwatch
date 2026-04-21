package daemon

import (
	"sync"
	"testing"
)

func TestSamplerConcurrentRecord(t *testing.T) {
	s := newSampler()
	const goroutines = 50
	const recsEach = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < recsEach; j++ {
				s.record(float64(j))
			}
		}()
	}
	wg.Wait()

	st := s.stats()
	expected := int64(goroutines * recsEach)
	if st.Count != expected {
		t.Fatalf("expected count %d, got %d", expected, st.Count)
	}
}

func TestSamplerConcurrentResetAndRecord(t *testing.T) {
	s := newSampler()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			s.record(float64(i))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			s.reset()
		}
	}()

	wg.Wait()
	// Just ensure no race or panic — stats must be consistent.
	st := s.stats()
	if st.Count < 0 {
		t.Fatalf("count must not be negative")
	}
}
