// Package notifier delivers port-change events to configured output channels.
//
// # Overview
//
// A Notifier wraps one or more delivery mechanisms:
//
//   - Standard output: every event is printed in a human-readable line.
//   - Webhook (optional): events are serialised as JSON and POSTed to a
//     user-supplied URL, enabling integration with alerting systems such as
//     PagerDuty, Slack incoming webhooks, or custom HTTP receivers.
//
// # Usage
//
//	n := notifier.New("https://hooks.example.com/portwatch", true)
//	err := n.Notify(notifier.Event{
//		Type:     "opened",
//		Port:     8080,
//		Protocol: "tcp",
//	})
//
// Passing an empty string as the webhook URL disables HTTP delivery; only
// stdout output is produced in that case.
//
// # Event fields
//
//   - Timestamp: set automatically to time.Now() when zero.
//   - Type: either "opened" or "closed".
//   - Port: the affected port number.
//   - Protocol: "tcp" or "udp".
package notifier
