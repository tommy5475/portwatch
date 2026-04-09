package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.ScanInterval != 30*time.Second {
		t.Errorf("expected scan interval 30s, got %v", cfg.ScanInterval)
	}

	if len(cfg.PortRanges) != 1 {
		t.Errorf("expected 1 port range, got %d", len(cfg.PortRanges))
	}

	if !cfg.AlertConfig.Enabled {
		t.Error("expected alerts to be enabled by default")
	}

	if cfg.LogLevel != "info" {
		t.Errorf("expected log level 'info', got %s", cfg.LogLevel)
	}
}

func TestLoadFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configJSON := `{
		"scan_interval": "60s",
		"port_ranges": [
			{"start": 80, "end": 443},
			{"start": 8000, "end": 8080}
		],
		"alert_config": {
			"enabled": true,
			"batch_size": 10,
			"channels": ["stdout", "email"]
		},
		"log_level": "debug"
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.ScanInterval != 60*time.Second {
		t.Errorf("expected scan interval 60s, got %v", cfg.ScanInterval)
	}

	if len(cfg.PortRanges) != 2 {
		t.Errorf("expected 2 port ranges, got %d", len(cfg.PortRanges))
	}

	if cfg.AlertConfig.BatchSize != 10 {
		t.Errorf("expected batch size 10, got %d", cfg.AlertConfig.BatchSize)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  Default(),
			wantErr: false,
		},
		{
			name: "invalid scan interval",
			config: &Config{
				ScanInterval: 500 * time.Millisecond,
				PortRanges:   []PortRange{{Start: 1, End: 100}},
			},
			wantErr: true,
		},
		{
			name: "no port ranges",
			config: &Config{
				ScanInterval: 30 * time.Second,
				PortRanges:   []PortRange{},
			},
			wantErr: true,
		},
		{
			name: "invalid port range",
			config: &Config{
				ScanInterval: 30 * time.Second,
				PortRanges:   []PortRange{{Start: 100, End: 50}},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
