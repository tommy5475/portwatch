package daemon

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestHooksRegisterAndFire(t *testing.T) {
	h := newHooks()
	var called int32
	h.Register(HookAfterScan, func(_ hookEvent, _ any) {
		atomic.AddInt32(&called, 1)
	})
	h.Fire(HookAfterScan, nil)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected 1 call, got %d", called)
	}
}

func TestHooksMultipleCallbacks(t *testing.T) {
	h := newHooks()
	var count int32
	for i := 0; i < 3; i++ {
		h.Register(HookBeforeScan, func(_ hookEvent, _ any) {
			atomic.AddInt32(&count, 1)
		})
	}
	h.Fire(HookBeforeScan, nil)
	if atomic.LoadInt32(&count) != 3 {
		t.Fatalf("expected 3 calls, got %d", count)
	}
}

func TestHooksEventIsolation(t *testing.T) {
	h := newHooks()
	var fired bool
	h.Register(HookOnChange, func(_ hookEvent, _ any) { fired = true })
	h.Fire(HookOnDegraded, nil) // different event — should not trigger
	if fired {
		t.Fatal("hook fired for wrong event")
	}
}

func TestHooksClear(t *testing.T) {
	h := newHooks()
	var count int32
	h.Register(HookOnRecovered, func(_ hookEvent, _ any) { atomic.AddInt32(&count, 1) })
	h.Clear(HookOnRecovered)
	h.Fire(HookOnRecovered, nil)
	if atomic.LoadInt32(&count) != 0 {
		t.Fatal("hook fired after Clear")
	}
}

func TestHooksLen(t *testing.T) {
	h := newHooks()
	if h.Len(HookAfterScan) != 0 {
		t.Fatal("expected 0 before registration")
	}
	h.Register(HookAfterScan, func(_ hookEvent, _ any) {})
	h.Register(HookAfterScan, func(_ hookEvent, _ any) {})
	if h.Len(HookAfterScan) != 2 {
		t.Fatalf("expected 2, got %d", h.Len(HookAfterScan))
	}
}

func TestHooksNilFnIgnored(t *testing.T) {
	h := newHooks()
	h.Register(HookBeforeScan, nil) // must not panic
	if h.Len(HookBeforeScan) != 0 {
		t.Fatal("nil fn should not be registered")
	}
}

func TestHooksConcurrentFireRegister(t *testing.T) {
	h := newHooks()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.Register(HookAfterScan, func(_ hookEvent, _ any) {})
			h.Fire(HookAfterScan, nil)
		}()
	}
	wg.Wait()
}
