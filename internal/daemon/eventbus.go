package daemon

import "sync"

// eventBus is a simple publish/subscribe bus for internal daemon events.
// Subscribers register a channel and receive copies of every published event.
// All methods are safe for concurrent use.
type eventBus struct {
	mu   sync.RWMutex
	subs map[string][]chan Event
}

// Event carries a topic and an arbitrary payload.
type Event struct {
	Topic   string
	Payload interface{}
}

func newEventBus() *eventBus {
	return &eventBus{subs: make(map[string][]chan Event)}
}

// Subscribe returns a buffered channel that will receive events for topic.
// The caller must drain or close the channel before calling Unsubscribe.
func (b *eventBus) Subscribe(topic string, buf int) chan Event {
	if buf <= 0 {
		buf = 8
	}
	ch := make(chan Event, buf)
	b.mu.Lock()
	b.subs[topic] = append(b.subs[topic], ch)
	b.mu.Unlock()
	return ch
}

// Unsubscribe removes ch from the subscriber list for topic and closes it.
func (b *eventBus) Unsubscribe(topic string, ch chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	list := b.subs[topic]
	for i, s := range list {
		if s == ch {
			b.subs[topic] = append(list[:i], list[i+1:]...)
			close(ch)
			return
		}
	}
}

// Publish sends ev to all subscribers of ev.Topic.
// Slow subscribers are skipped (non-blocking send).
func (b *eventBus) Publish(ev Event) {
	b.mu.RLock()
	list := b.subs[ev.Topic]
	b.mu.RUnlock()
	for _, ch := range list {
		select {
		case ch <- ev:
		default:
		}
	}
}

// Drain closes and removes all subscribers for every topic.
func (b *eventBus) Drain() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for topic, list := range b.subs {
		for _, ch := range list {
			close(ch)
		}
		delete(b.subs, topic)
	}
}
