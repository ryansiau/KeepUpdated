package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/ryansiau/utilities/go/config"
	graceful_shutdown "github.com/ryansiau/utilities/go/pkg/graceful-shutdown"
	"github.com/ryansiau/utilities/go/worker"
)

func main() {
	// Get config directory
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	logrus.Infof("Using config file: %s\n", configPath)

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logrus.Fatalf("Error loading config: %v\n", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logrus.Fatalf("Validation failed: %v\n", err)
	}

	logrus.Info("Configuration loaded.")

	// Create Worker
	sm := graceful_shutdown.NewGracefulShutdown()
	w, err := worker.NewWorker(cfg, sm)
	if err != nil {
		logrus.Fatalf("Error creating worker: %v\n", err)
	}

	// Start Worker
	if err = w.Run(); err != nil {
		logrus.Fatalf("Error running worker: %v\n", err)
	}
}
