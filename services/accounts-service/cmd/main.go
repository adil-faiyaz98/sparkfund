package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/adapters/grpc"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/adapters/postgres"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/application"
)

var (
	port  = flag.Int("port", 50051, "gRPC server port")
	dbDSN = flag.String("db-dsn", "postgres://postgres:postgres@localhost:5432/money_pulse?sslmode=disable", "Database DSN")
)

func main() {
	flag.Parse()

	// Initialize database adapter
	db, err := postgres.NewAdapter(*dbDSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize application core
	app := application.NewAccountService(db)

	// Initialize gRPC server
	server := grpc.NewServer(app)

	// Create context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start gRPC server in a goroutine
	go func() {
		if err := server.Run(*port); err != nil {
			log.Printf("Failed to run gRPC server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	// Perform graceful shutdown
	fmt.Println("Shutting down gracefully...")
}
