package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	traefikkeymate "github.com/numkem/traefik-keymate"
	"github.com/numkem/traefik-keymate/keymate"
	log "github.com/sirupsen/logrus"
)

func main() {
	configFilename := flag.String("config", "", "configuration filename")
	flag.Parse()

	if *configFilename == "" {
		log.Fatal("configuration filename required")
	}

	log.Info("Traefik Keymate starting...")

	cfg, err := traefikkeymate.NewConfig(*configFilename)
	if err != nil {
		log.Fatalf("Failed to read configuraiton: %v", err)
	}

	// Create manager connection
	mgr, err := keymate.NewEtcdManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	errs := mgr.ApplyConfig(ctx, cfg)
	for _, err := range errs {
		log.Errorf("Error found while applying configuration: %v", err)
	}

	for _, err := range errs {
		log.Error(err)
	}

	log.Info("Configuration applied. Exiting")
}
