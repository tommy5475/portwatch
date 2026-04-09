// Package alert provides alerting functionality for port changes.
package alert

import (
	"fmt"
	"log"
	"os"
	"time"
)

// PortChange represents a change in port state.
type PortChange struct {
	Port      int
	OldState  bool
	NewState  bool
	Timestamp time.Time
}

// Alerter handles alerting for port changes.
type Alerter struct {
	logger *log.Logger
	output *os.File
}

// New creates a new Alerter instance.
func New(output *os.File) *Alerter {
	if output == nil {
		output = os.Stdout
	}
	return &Alerter{
		logger: log.New(output, "[PORTWATCH] ", log.LstdFlags),
		output: output,
	}
}

// Alert sends an alert for a port change.
func (a *Alerter) Alert(change PortChange) error {
	if a == nil {
		return fmt.Errorf("alerter is nil")
	}

	var msg string
	if change.NewState {
		msg = fmt.Sprintf("PORT OPENED: Port %d is now OPEN (was CLOSED)", change.Port)
	} else {
		msg = fmt.Sprintf("PORT CLOSED: Port %d is now CLOSED (was OPEN)", change.Port)
	}

	a.logger.Println(msg)
	return nil
}

// AlertBatch sends alerts for multiple port changes.
func (a *Alerter) AlertBatch(changes []PortChange) error {
	if a == nil {
		return fmt.Errorf("alerter is nil")
	}

	if len(changes) == 0 {
		return nil
	}

	a.logger.Printf("Detected %d port change(s):", len(changes))
	for _, change := range changes {
		if err := a.Alert(change); err != nil {
			return err
		}
	}

	return nil
}

// Info logs an informational message.
func (a *Alerter) Info(msg string) {
	if a != nil {
		a.logger.Println("INFO:", msg)
	}
}

// Error logs an error message.
func (a *Alerter) Error(msg string) {
	if a != nil {
		a.logger.Println("ERROR:", msg)
	}
}
