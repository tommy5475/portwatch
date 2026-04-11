package daemon

import (
	"sync"
	"testing"
	"time"
)

func TestFanoutNoSubscribers(t *testing.T) {
	f := newFanout[int](4)
	// publish with no subscribers must not block or panic
	f.publish(42)
}

func TestFanoutSingleSubscriberReceives(t *testing.T) {
	f := newFanout[int](4)
	ch, cancel := f.subscribe()
	defer cancel()

	f.publish(7)

	select {
	case v := <-ch:
		if v != 7 {
			t.Fatalf("expected 7, got %d", v)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for value")
	}
}

func TestFanoutMultipleSubscribersAllReceive(t *testing.T) {
	f := newFanout[string](4)
	const n = 5
	chans := make([]<-chan string, n)
	cancels := make([]func(), n)
	for i := range n {
		chans[i], cancels[i] = f.subscribe()
		defer cancels[i]()
	}

	f.publish("hello")

	for i, ch := range chans {
		select {
		case v := <-ch:
			if v != "hello" {
				t.Fatalf("sub %d: expected hello, got %s", i, v)
			}
		case <-time.After(time.Second):
			t.Fatalf("sub %d: timed out", i)
		}
	}
}

func TestFanoutCancelRemovesSubscriber(t *testing.T) {
	f := newFanout[int](4)
	_, cancel := f.subscribe()
	if f.len() != 1 {
		t.Fatalf("expected 1 subscriber, got %d", f.len())
	}
	cancel()
	if f.len() != 0 {
		t.Fatalf("expected 0 subscribers after cancel, got %d", f.len())
	}
}

func TestFanoutSlowConsumerSkipped(t *testing.T) {
	f := newFanout[int](1)
	slow, cancelSlow := f.subscribe()
	fast, cancelFast := f.subscribe()
	defer cancelSlow()
	defer cancelFast()

	// Fill slow subscriber's buffer so the second publish is dropped for it.
	f.publish(1)
	f.publish(2) // slow buffer full; fast should still get both

	// fast should receive both values
	for range 2 {
		select {
		case <-fast:
		case <-time.After(time.Second):
			t.Fatal("fast subscriber timed out")
		}
	}
	_ = slow
}

func TestFanoutConcurrentPublishSubscribe(t *testing.T) {
	f := newFanout[int](8)
	var wg sync.WaitGroup
	const publishers = 4
	const msgs = 20

	ch, cancel := f.subscribe()
	defer cancel()

	for range publishers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range msgs {
				f.publish(i)
			}
		}()
	}

	wg.Wait()
	// drain without blocking
	done := time.After(500 * time.Millisecond)
	for {
		select {
		case <-ch:
		case <-done:
			return
		}
	}
}
