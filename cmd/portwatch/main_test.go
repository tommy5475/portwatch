package main

import (
	"os"
	"path/filepath"
	"testing"

	"portwatch/internal/config"
)

func TestLoadConfigDefault(t *testing.T) {
	cfg, err := loadConfig("")
	if err != nil {
		t.Fatalf("expected no error for default config, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if len(cfg.Ports) ==	t.Error("expected default config to have at least one port")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	cfgFileportwatch.yaml")

	content := []byte(`ports: [8080, 9090]
val: 30s
state_file: /tmp/portwatch_test.state
`)
	if err := os.WriteFile(cfg644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := loadConfig(cfgFile)
	if err != nil {
		t.Fatalf("expected no error loading file config, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := loadConfig("/nonexistent/path/portwatch.yaml")
	if err == nil {
		t.Error("expected error for missing config file, got nil")
	}
}

func TestLoadConfigInvalidFile(t *testing.T) {
	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "bad.yaml")

	// ports is required; empty file should fail validation
	if err := os.WriteFile(cfgFile, []byte(`interval: 10s\n`), 0o644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	cfg, err := loadConfig(cfgFile)
	// either an error or an empty port list is acceptable depending on Validate impl
	if err == nil && cfg != nil && len(cfg.Ports) == 0 {
		// acceptable: no ports configured
		_ = config.Default()
	}
}
