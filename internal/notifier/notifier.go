// Package notifier provides output channels for port change alerts,
// supporting stdout and webhook-based notifications.
package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Event represents a port change notification payload.
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "opened" or "closed"
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
}

// Notifier dispatches events to one or more outputs.
type Notifier struct {
	webhookURL string
	client     *http.Client
	verbose    bool
}

// New creates a Notifier. webhookURL may be empty to disable webhook output.
func New(webhookURL string, verbose bool) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 5 * time.Second},
		verbose:    verbose,
	}
}

// Notify prints the event to stdout and, if configured, sends it to the webhook.
func (n *Notifier) Notify(e Event) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}

	fmt.Printf("[%s] port %s/%d %s\n",
		e.Timestamp.Format(time.RFC3339), e.Protocol, e.Port, e.Type)

	if n.webhookURL == "" {
		return nil
	}
	return n.sendWebhook(e)
}

func (n *Notifier) sendWebhook(e Event) error {
	body, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("notifier: marshal event: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notifier: webhook post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("notifier: webhook returned status %d", resp.StatusCode)
	}
	if n.verbose {
		fmt.Printf("[notifier] webhook delivered (status %d)\n", resp.StatusCode)
	}
	return nil
}
