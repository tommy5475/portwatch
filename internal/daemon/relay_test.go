package daemon

import (
	"testing"
	"time"
)

func TestRelayNoSubscribers(t *testing.T) {
	r := newRelay[int](4, nil)
	r.send(42) // must not panic
}

func TestRelaySubscriberReceives(t *testing.T) {
	r := newRelay[int](4, nil)
	ch, cancel := r.subscribe()
	defer cancel()
	r.send(7)
	select {
	case v := <-ch:
		if v != 7 {
			t.Fatalf("want 7, got %d", v)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for value")
	}
}

func TestRelayFilterDropsEvents(t *testing.T) {
	r := newRelay[int](4, func(v int) bool { return v > 10 })
	ch, cancel := r.subscribe()
	defer cancel()
	r.send(5)
	r.send(20)
	select {
	case v := <-ch:
		if v != 20 {
			t.Fatalf("want 20, got %d", v)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
	if len(ch) != 0 {
		t.Fatal("unexpected extra message")
	}
}

func TestRelayCancelRemovesSubscriber(t *testing.T) {
	r := newRelay[int](4, nil)
	_, cancel := r.subscribe()
	if r.len() != 1 {
		t.Fatalf("want 1 subscriber, got %d", r.len())
	}
	cancel()
	if r.len() != 0 {
		t.Fatalf("want 0 subscribers after cancel, got %d", r.len())
	}
}

func TestRelayMultipleSubscribersAllReceive(t *testing.T) {
	r := newRelay[string](4, nil)
	ch1, c1 := r.subscribe()
	ch2, c2 := r.subscribe()
	defer c1()
	defer c2()
	r.send("hello")
	for _, ch := range []<-chan string{ch1, ch2} {
		select {
		case v := <-ch:
			if v != "hello" {
				t.Fatalf("want hello, got %s", v)
			}
		case <-time.After(time.Second):
			t.Fatal("timeout")
		}
	}
}

func TestRelayDefaultBufSize(t *testing.T) {
	r := newRelay[int](0, nil)
	if r.bufSize != 16 {
		t.Fatalf("want default bufSize 16, got %d", r.bufSize)
	}
}
