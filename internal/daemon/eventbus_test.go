package daemon

import (
	"testing"
	"time"
)

func TestEventBusSubscribeReceives(t *testing.T) {
	bus := newEventBus()
	ch := bus.Subscribe("scan.done", 4)
	bus.Publish(Event{Topic: "scan.done", Payload: 42})
	select {
	case ev := <-ch:
		if ev.Payload.(int) != 42 {
			t.Fatalf("expected 42, got %v", ev.Payload)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestEventBusTopicIsolation(t *testing.T) {
	bus := newEventBus()
	ch := bus.Subscribe("scan.done", 4)
	bus.Publish(Event{Topic: "scan.error", Payload: "err"})
	select {
	case <-ch:
		t.Fatal("should not receive event for different topic")
	case <-time.After(20 * time.Millisecond):
	}
}

func TestEventBusUnsubscribeClosesChannel(t *testing.T) {
	bus := newEventBus()
	ch := bus.Subscribe("change.found", 4)
	bus.Unsubscribe("change.found", ch)
	_, open := <-ch
	if open {
		t.Fatal("channel should be closed after Unsubscribe")
	}
}

func TestEventBusSlowConsumerDropped(t *testing.T) {
	bus := newEventBus()
	// buffer of 1 – second publish must not block
	bus.Subscribe("alert.sent", 1)
	bus.Publish(Event{Topic: "alert.sent", Payload: "a"})
	bus.Publish(Event{Topic: "alert.sent", Payload: "b"}) // must not block
}

func TestEventBusMultipleSubscribers(t *testing.T) {
	bus := newEventBus()
	ch1 := bus.Subscribe("scan.done", 4)
	ch2 := bus.Subscribe("scan.done", 4)
	bus.Publish(Event{Topic: "scan.done", Payload: "x"})
	for _, ch := range []chan Event{ch1, ch2} {
		select {
		case ev := <-ch:
			if ev.Payload.(string) != "x" {
				t.Fatalf("unexpected payload %v", ev.Payload)
			}
		case <-time.After(time.Second):
			t.Fatal("timed out")
		}
	}
}

func TestEventBusDrainClosesAll(t *testing.T) {
	bus := newEventBus()
	ch1 := bus.Subscribe("scan.done", 4)
	ch2 := bus.Subscribe("scan.error", 4)
	bus.Drain()
	for _, ch := range []chan Event{ch1, ch2} {
		_, open := <-ch
		if open {
			t.Fatal("channel should be closed after Drain")
		}
	}
}

func TestEventBusDefaultBufferOnZero(t *testing.T) {
	bus := newEventBus()
	ch := bus.Subscribe("scan.done", 0)
	// publish up to default buffer (8) without blocking
	for i := 0; i < 8; i++ {
		bus.Publish(Event{Topic: "scan.done", Payload: i})
	}
	if len(ch) != 8 {
		t.Fatalf("expected 8 buffered events, got %d", len(ch))
	}
}
