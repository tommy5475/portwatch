package notifier_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portwatch/internal/notifier"
)

func TestNotifyStdoutOnly(t *testing.T) {
	n := notifier.New("", false)
	e := notifier.Event{
		Timestamp: time.Now(),
		Type:      "opened",
		Port:      8080,
		Protocol:  "tcp",
	}
	if err := n.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNotifyWebhook(t *testing.T) {
	var received notifier.Event

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notifier.New(srv.URL, true)
	e := notifier.Event{
		Timestamp: time.Now(),
		Type:      "closed",
		Port:      443,
		Protocol:  "tcp",
	}
	if err := n.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Port != e.Port {
		t.Errorf("expected port %d, got %d", e.Port, received.Port)
	}
	if received.Type != e.Type {
		t.Errorf("expected type %q, got %q", e.Type, received.Type)
	}
}

func TestNotifyWebhookError(t *testing.T) {
	n := notifier.New("http://127.0.0.1:0/webhook", false)
	e := notifier.Event{Type: "opened", Port: 22, Protocol: "tcp"}
	if err := n.Notify(e); err == nil {
		t.Fatal("expected error for unreachable webhook, got nil")
	}
}

func TestNotifyZeroTimestamp(t *testing.T) {
	var received notifier.Event
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	n := notifier.New(srv.URL, false)
	// zero timestamp should be filled in automatically
	if err := n.Notify(notifier.Event{Type: "opened", Port: 80, Protocol: "tcp"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp in webhook payload")
	}
}
