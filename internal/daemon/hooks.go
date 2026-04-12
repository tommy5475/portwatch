package daemon

import (
	"sync"
)

// hookEvent represents a lifecycle event that can trigger registered hooks.
type hookEvent string

const (
	HookBeforeScan  hookEvent = "before_scan"
	HookAfterScan   hookEvent = "after_scan"
	HookOnChange    hookEvent = "on_change"
	HookOnDegraded  hookEvent = "on_degraded"
	HookOnRecovered hookEvent = "on_recovered"
)

// HookFn is a function invoked when a lifecycle event fires.
type HookFn func(event hookEvent, payload any)

// hooks manages a registry of lifecycle callbacks keyed by event type.
type hooks struct {
	mu       sync.RWMutex
	registry map[hookEvent][]HookFn
}

func newHooks() *hooks {
	return &hooks{
		registry: make(map[hookEvent][]HookFn),
	}
}

// Register adds fn to the list of callbacks for the given event.
// Multiple functions may be registered for the same event.
func (h *hooks) Register(event hookEvent, fn HookFn) {
	if fn == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.registry[event] = append(h.registry[event], fn)
}

// Fire invokes all callbacks registered for event, passing payload.
// Callbacks are called synchronously in registration order.
func (h *hooks) Fire(event hookEvent, payload any) {
	h.mu.RLock()
	fns := make([]HookFn, len(h.registry[event]))
	copy(fns, h.registry[event])
	h.mu.RUnlock()

	for _, fn := range fns {
		fn(event, payload)
	}
}

// Clear removes all callbacks registered for the given event.
func (h *hooks) Clear(event hookEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.registry, event)
}

// Len returns the number of callbacks registered for event.
func (h *hooks) Len(event hookEvent) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.registry[event])
}
