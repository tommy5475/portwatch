package daemon

import (
	"sync/atomic"
	"testing"
	"time"
)

// TestHooksPayloadPropagation verifies that the payload passed to Fire
// reaches every registered callback unchanged.
func TestHooksPayloadPropagation(t *testing.T) {
	h := newHooks()
	type scanResult struct{ duration time.Duration }
	want := scanResult{duration: 42 * time.Millisecond}

	var got scanResult
	h.Register(HookAfterScan, func(_ hookEvent, p any) {
		got = p.(scanResult)
	})
	h.Fire(HookAfterScan, want)

	if got != want {
		t.Fatalf("payload mismatch: got %v, want %v", got, want)
	}
}

// TestHooksFireUnregisteredEventIsNoop ensures firing an event with no
// registered callbacks does not panic.
func TestHooksFireUnregisteredEventIsNoop(t *testing.T) {
	h := newHooks()
	h.Fire(hookEvent("nonexistent"), nil) // must not panic
}

// TestHooksDegradedRecoveredCycle simulates the watcher degraded/recovered
// lifecycle and asserts hooks fire in the correct order.
func TestHooksDegradedRecoveredCycle(t *testing.T) {
	h := newHooks()
	var seq []hookEvent

	record := func(e hookEvent, _ any) { seq = append(seq, e) }
	h.Register(HookOnDegraded, record)
	h.Register(HookOnRecovered, record)

	h.Fire(HookOnDegraded, nil)
	h.Fire(HookOnRecovered, nil)

	if len(seq) != 2 || seq[0] != HookOnDegraded || seq[1] != HookOnRecovered {
		t.Fatalf("unexpected sequence: %v", seq)
	}
}

// TestHooksClearDoesNotAffectOtherEvents ensures Clear is scoped to a
// single event type.
func TestHooksClearDoesNotAffectOtherEvents(t *testing.T) {
	h := newHooks()
	var fired int32
	h.Register(HookOnChange, func(_ hookEvent, _ any) { atomic.AddInt32(&fired, 1) })
	h.Register(HookAfterScan, func(_ hookEvent, _ any) { atomic.AddInt32(&fired, 1) })

	h.Clear(HookOnChange)
	h.Fire(HookOnChange, nil)  // cleared — should not fire
	h.Fire(HookAfterScan, nil) // still registered

	if atomic.LoadInt32(&fired) != 1 {
		t.Fatalf("expected 1 call, got %d", fired)
	}
}
