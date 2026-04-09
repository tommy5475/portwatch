package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/monitor"
	"portwatch/internal/state"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	cfgPath := flag.String("config", "", "path to config file (optional)")
	showVersion := flag.Bool("version", false, "printlag.Parse()

n		fmt.Printf("portwatch %s (%s)\n", version, commit)
		os.Exit(0)
	}

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	st, err := state.New(cfg.StateFile)
	if err != nil {
		log.Fatalf("state error: %v", err)
	}

	al := alert.New(cfg)
	mon := monitor.New(cfg, st)

	if err := mon.Start(); err != nil {
		log.Fatalf("monitor start error: %v", err)
	}
	defer mon.Stop()

	log.Printf("portwatch started — watching ports %v", cfg.Ports)

	go func() {
		for changes := range mon.Changes() {
			if len(changes) > 0 {
				if err := al.AlertBatch(changes); err != nil {
					log.Printf("alert error: %v", err)
				}
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("portwatch shutting down")
}

func loadConfig(path string) (*config.Config, error) {
	if path != "" {
		return config.LoadFromFile(path)
	}
	cfg := config.Default()
	return cfg, cfg.Validate()
}
