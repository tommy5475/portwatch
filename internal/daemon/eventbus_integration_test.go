package daemon

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestEventBusConcurrentPublishSubscribe hammers the bus from many goroutines
// to surface data races (run with -race).
func TestEventBusConcurrentPublishSubscribe(t *testing.T) {
	bus := newEventBus()
	const workers = 20
	const msgs = 50

	var received int64
	var wg sync.WaitGroup

	// start subscribers
	for i := 0; i < workers/2; i++ {
		ch := bus.Subscribe("scan.done", 64)
		wg.Add(1)
		go func(c chan Event) {
			defer wg.Done()
			for range c {
				atomic.AddInt64(&received, 1)
			}
		}(ch)
	}

	// start publishers
	var pub sync.WaitGroup
	for i := 0; i < workers/2; i++ {
		pub.Add(1)
		go func() {
			defer pub.Done()
			for j := 0; j < msgs; j++ {
				bus.Publish(Event{Topic: "scan.done", Payload: j})
			}
		}()
	}

	pub.Wait()
	// give subscribers a moment to drain buffered events
	time.Sleep(50 * time.Millisecond)
	bus.Drain()
	wg.Wait()

	t.Logf("received %d events across %d subscribers", received, workers/2)
}
