package monitor

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ports := []int{80, 443, 8080}
	interval := 5 * time.Second

	m := New(ports, interval)

	if m == nil {
		t.Fatal("New() returned nil")
	}

	if len(m.ports) != len(ports) {
		t.Errorf("expected %d ports, got %d", len(ports), len(m.ports))
	}

	if m.interval != interval {
		t.Errorf("expected interval %v, got %v", interval, m.interval)
	}

	if m.lastState == nil {
		t.Error("lastState map not initialized")
	}

	if m.changesChan == nil {
		t.Error("changesChan not initialized")
	}
}

func TestMonitorStart(t *testing.T) {
	ports := []int{9999} // Use unlikely port
	interval := 100 * time.Millisecond

	m := New(ports, interval)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	err := m.Start(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}

	// Verify initial state was captured
	state := m.GetState()
	if _, exists := state[9999]; !exists {
		t.Error("expected port 9999 in state after monitoring")
	}
}

func TestGetState(t *testing.T) {
	ports := []int{80, 443}
	interval := 1 * time.Second

	m := New(ports, interval)

	// Manually set some state
	m.mu.Lock()
	m.lastState[80] = true
	m.lastState[443] = false
	m.mu.Unlock()

	state := m.GetState()

	if len(state) != 2 {
		t.Errorf("expected 2 ports in state, got %d", len(state))
	}

	if state[80] != true {
		t.Error("expected port 80 to be open")
	}

	if state[443] != false {
		t.Error("expected port 443 to be closed")
	}
}

func TestChangesChannel(t *testing.T) {
	ports := []int{8080}
	interval := 1 * time.Second

	m := New(ports, interval)

	changesChan := m.Changes()
	if changesChan == nil {
		t.Fatal("Changes() returned nil channel")
	}

	// Verify we can receive from the channel
	go func() {
		m.changesChan <- PortChange{
			Port:     8080,
			WasOpen:  false,
			NowOpen:  true,
			Detected: time.Now(),
		}
	}()

	select {
	case change := <-changesChan:
		if change.Port != 8080 {
			t.Errorf("expected port 8080, got %d", change.Port)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for change")
	}
}
