package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adil-faiyaz98/money-pulse/pkg/swagger"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/repository"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/service"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/transport/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[accounts-service] ", log.LstdFlags)

	// Initialize repository (implementation will depend on your database choice)
	// For now, we'll use a mock repository
	repo := repository.NewMockAccountRepository()

	// Initialize service
	accountService := service.NewAccountService(repo)

	// Initialize HTTP handler
	handler := http.NewAccountHandler(accountService)

	// Initialize Gin router
	router := gin.Default()

	// Serve static files for Swagger UI
	router.Static("/swagger", "./api/accounts/v1")
	router.GET("/docs", swagger.SwaggerUI("/swagger/openapi.yaml"))

	// Register routes
	handler.RegisterRoutes(router)

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Create shutdown context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exiting")
}
