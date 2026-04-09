// Package config provides configuration management for portwatch.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config represents the portwatch configuration.
type Config struct {
	// ScanInterval is the duration between scans
	ScanInterval time.Duration `json:"scan_interval"`
	// PortRanges defines the port ranges to monitor
	PortRanges []PortRange `json:"port_ranges"`
	// AlertConfig contains alert configuration
	AlertConfig AlertConfig `json:"alert_config"`
	// LogLevel defines the logging verbosity
	LogLevel string `json:"log_level"`
}

// PortRange defines a range of ports to scan.
type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// AlertConfig contains alert-specific configuration.
type AlertConfig struct {
	// Enabled determines if alerts are active
	Enabled bool `json:"enabled"`
	// BatchSize is the number of changes to batch before alerting
	BatchSize int `json:"batch_size"`
	// Channels defines which alert channels to use
	Channels []string `json:"channels"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		ScanInterval: 30 * time.Second,
		PortRanges: []PortRange{
			{Start: 1, End: 1024},
		},
		AlertConfig: AlertConfig{
			Enabled:   true,
			BatchSize: 5,
			Channels:  []string{"stdout"},
		},
		LogLevel: "info",
	}
}

// LoadFromFile loads configuration from a JSON file.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.ScanInterval < time.Second {
		return fmt.Errorf("scan_interval must be at least 1 second")
	}

	if len(c.PortRanges) == 0 {
		return fmt.Errorf("at least one port range must be specified")
	}

	for _, pr := range c.PortRanges {
		if pr.Start < 1 || pr.Start > 65535 {
			return fmt.Errorf("invalid port range start: %d", pr.Start)
		}
		if pr.End < 1 || pr.End > 65535 {
			return fmt.Errorf("invalid port range end: %d", pr.End)
		}
		if pr.Start > pr.End {
			return fmt.Errorf("port range start must be <= end")
		}
	}

	return nil
}
