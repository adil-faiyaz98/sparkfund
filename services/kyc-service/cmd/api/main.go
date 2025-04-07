package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sparkfund/services/kyc-service/internal/app"
	"sparkfund/services/kyc-service/internal/config"
)

var (
	version    = "dev"
	commitHash = "none"
	buildTime  = "unknown"
)

func main() {
	// Initialize config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize application
	application, err := app.New(app.Config{
		Version:    version,
		CommitHash: commitHash,
		BuildTime:  buildTime,
		Config:     cfg,
	})
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start the application
	go func() {
		if err := application.Start(); err != nil {
			log.Fatalf("Failed to start application: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown application: %v", err)
	}
}
