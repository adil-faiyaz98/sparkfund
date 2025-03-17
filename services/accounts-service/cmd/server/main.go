package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/config"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/grpc/handler"
	httpHandler "github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/http"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/repository/postgres"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize repository
	repo, err := postgres.NewAccountRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to create account repository: %v", err)
	}
	defer repo.Close()

	// Initialize service
	accountService := service.NewAccountService(repo)

	// Start gRPC server
	grpcServer := grpc.NewServer()
	accountHandler := handler.NewAccountHandler(accountService)
	// Register gRPC service handler
	// pb.RegisterAccountServiceServer(grpcServer, accountHandler)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	go func() {
		log.Printf("Starting gRPC server on port %d", cfg.GrpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server with Gin
	router := gin.Default()
	httpHandlers := httpHandler.NewAccountHandler(accountService)
	httpHandlers.RegisterRoutes(router)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpPort),
		Handler: router,
	}

	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.HttpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()
	log.Println("Servers exited properly")
}
