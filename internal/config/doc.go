/*
Package config provides configuration management for portwatch.

The config package handles loading, validating, and managing the portwatch
configuration. Configuration can be loaded from JSON files or created with
sensible defaults.

Configuration Structure:

The main Config struct contains:
  - ScanInterval: How often to scan ports (minimum 1 second)
  - PortRanges: One or more port ranges to monitor
  - AlertConfig: Alert-specific settings including channels and batch size
  - LogLevel: Logging verbosity level

Usage:

	// Use default configuration
	cfg := config.Default()

	// Load from file
	cfg, err := config.LoadFromFile("/etc/portwatch/config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

Example Configuration File:

	{
	  "scan_interval": "30s",
	  "port_ranges": [
	    {"start": 1, "end": 1024},
	    {"start": 8000, "end": 9000}
	  ],
	  "alert_config": {
	    "enabled": true,
	    "batch_size": 5,
	    "channels": ["stdout", "email"]
	  },
	  "log_level": "info"
	}

Port Ranges:

Port ranges must be valid (1-65535) and start must be less than or equal to end.
Multiple ranges can be specified to monitor different port groups.

Alert Configuration:

The AlertConfig controls how alerts are generated and delivered:
  - Enabled: Master switch for alerts
  - BatchSize: Number of changes to accumulate before alerting
  - Channels: List of alert delivery methods (stdout, email, webhook, etc.)
*/
package config
