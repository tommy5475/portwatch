package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// PortState represents the state of a port at a specific time
type PortState struct {
	Port      int
	Open      bool
	Timestamp time.Time
}

// PortChange represents a change in port state
type PortChange struct {
	Port     int
	WasOpen  bool
	NowOpen  bool
	Detected time.Time
}

// Monitor watches ports for changes
type Monitor struct {
	scanner      *scanner.Scanner
	ports        []int
	interval     time.Duration
	mu           sync.RWMutex
	lastState    map[int]bool
	changesChan  chan PortChange
}

// New creates a new Monitor instance
func New(ports []int, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:     scanner.New(),
		ports:       ports,
		interval:    interval,
		lastState:   make(map[int]bool),
		changesChan: make(chan PortChange, 100),
	}
}

// Start begins monitoring ports
func (m *Monitor) Start(ctx context.Context) error {
	// Initial scan to establish baseline
	if err := m.scan(); err != nil {
		return fmt.Errorf("initial scan failed: %w", err)
	}

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(m.changesChan)
			return ctx.Err()
		case <-ticker.C:
			if err := m.scan(); err != nil {
				// Log error but continue monitoring
				continue
			}
		}
	}
}

// scan performs a port scan and detects changes
func (m *Monitor) scan() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, port := range m.ports {
		open, err := m.scanner.ScanTCPPort(port)
		if err != nil {
			return err
		}

		if lastOpen, exists := m.lastState[port]; exists && lastOpen != open {
			// State changed
			m.changesChan <- PortChange{
				Port:     port,
				WasOpen:  lastOpen,
				NowOpen:  open,
				Detected: time.Now(),
			}
		}

		m.lastState[port] = open
	}

	return nil
}

// Changes returns the channel for receiving port changes
func (m *Monitor) Changes() <-chan PortChange {
	return m.changesChan
}

// GetState returns the current state of all monitored ports
func (m *Monitor) GetState() map[int]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	state := make(map[int]bool, len(m.lastState))
	for k, v := range m.lastState {
		state[k] = v
	}
	return state
}
