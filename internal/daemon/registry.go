package daemon

import (
	"fmt"
	"sync"
)

// registry is a thread-safe map of named values with optional metadata.
// It is used to track named components, capabilities, or runtime entries
// that need to be discovered or iterated at runtime.
type registry struct {
	mu      sync.RWMutex
	entries map[string]registryEntry
}

type registryEntry struct {
	value interface{}
	tags  []string
}

func newRegistry() *registry {
	return &registry{
		entries: make(map[string]registryEntry),
	}
}

// Register adds or replaces a named entry with optional tags.
func (r *registry) Register(name string, value interface{}, tags ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[name] = registryEntry{value: value, tags: tags}
}

// Unregister removes a named entry. Returns true if it existed.
func (r *registry) Unregister(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.entries[name]
	if ok {
		delete(r.entries, name)
	}
	return ok
}

// Get returns the value for a name, or an error if not found.
func (r *registry) Get(name string) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[name]
	if !ok {
		return nil, fmt.Errorf("registry: %q not found", name)
	}
	return e.value, nil
}

// Has returns true if the name is registered.
func (r *registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.entries[name]
	return ok
}

// Names returns all registered names.
func (r *registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.entries))
	for k := range r.entries {
		names = append(names, k)
	}
	return names
}

// Len returns the number of registered entries.
func (r *registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}

// FilterByTag returns names whose tags include the given tag.
func (r *registry) FilterByTag(tag string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []string
	for name, e := range r.entries {
		for _, t := range e.tags {
			if t == tag {
				result = append(result, name)
				break
			}
		}
	}
	return result
}
