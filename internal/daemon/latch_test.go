package daemon

import (
	"sync"
	"testing"
)

func TestLatchInitiallyUnset(t *testing.T) {
	l := newLatch()
	if l.IsSet() {
		t.Fatal("expected latch to be unset initially")
	}
}

func TestLatchSetBecomesSet(t *testing.T) {
	l := newLatch()
	l.Set()
	if !l.IsSet() {
		t.Fatal("expected latch to be set after Set()")
	}
}

func TestLatchSetIsIdempotent(t *testing.T) {
	l := newLatch()
	l.Set()
	l.Set() // second call must not panic or change state
	if !l.IsSet() {
		t.Fatal("expected latch to remain set")
	}
}

func TestLatchIfSetCallsFnWhenSet(t *testing.T) {
	l := newLatch()
	l.Set()
	called := false
	result := l.IfSet(func() { called = true })
	if !called {
		t.Fatal("expected fn to be called")
	}
	if !result {
		t.Fatal("expected IfSet to return true")
	}
}

func TestLatchIfSetSkipsFnWhenUnset(t *testing.T) {
	l := newLatch()
	called := false
	result := l.IfSet(func() { called = true })
	if called {
		t.Fatal("expected fn not to be called")
	}
	if result {
		t.Fatal("expected IfSet to return false")
	}
}

func TestLatchConcurrentSet(t *testing.T) {
	l := newLatch()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Set()
			_ = l.IsSet()
		}()
	}
	wg.Wait()
	if !l.IsSet() {
		t.Fatal("expected latch to be set after concurrent Set calls")
	}
}
