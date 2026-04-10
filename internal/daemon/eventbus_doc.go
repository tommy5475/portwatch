// Package daemon contains the runtime components of portwatch.
//
// # EventBus
//
// eventBus provides a lightweight, in-process publish/subscribe mechanism
// used to decouple daemon subsystems.
//
// Typical usage:
//
//	bus := newEventBus()
//
//	// producer
//	bus.Publish(Event{Topic: "scan.done", Payload: result})
//
//	// consumer
//	ch := bus.Subscribe("scan.done", 16)
//	go func() {
//		for ev := range ch {
//			_ = ev.Payload
//		}
//	}()
//
// Defined topics used within the daemon:
//
//	"scan.done"    – emitted after every successful port scan
//	"scan.error"   – emitted when a scan returns an error
//	"change.found" – emitted when a state diff is non-empty
//	"alert.sent"   – emitted after an alert notification is dispatched
//
// Slow consumers are never blocked; events are dropped for full channels.
// Call Drain to shut down all subscriptions cleanly on daemon exit.
package daemon
