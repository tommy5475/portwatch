package daemon

// Well-known event topics published on the daemon's internal eventBus.
// Use these constants instead of raw strings to avoid typos and to make
// topic usage grep-able across the codebase.
const (
	// TopicScanDone is published after every successful port scan.
	// Payload: snapshot returned by newSnapshot.
	TopicScanDone = "scan.done"

	// TopicScanError is published when a scan attempt returns an error.
	// Payload: error value.
	TopicScanError = "scan.error"

	// TopicChangeFound is published when a state diff contains at least one
	// opened or closed port.
	// Payload: state.Diff value.
	TopicChangeFound = "change.found"

	// TopicAlertSent is published after a notification has been dispatched
	// successfully by the notifier.
	// Payload: alert.Alert value.
	TopicAlertSent = "alert.sent"

	// TopicHealthDegraded is published when the watcher transitions to the
	// degraded state.
	// Payload: string reason.
	TopicHealthDegraded = "health.degraded"

	// TopicHealthRecovered is published when the watcher returns to a healthy
	// state after being degraded.
	// Payload: nil.
	TopicHealthRecovered = "health.recovered"
)
