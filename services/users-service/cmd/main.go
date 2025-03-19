package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adil-faiyaz98/money-pulse/pkg/auth"
	"github.com/adil-faiyaz98/money-pulse/pkg/config"
	"github.com/adil-faiyaz98/money-pulse/pkg/logger"
	"github.com/adil-faiyaz98/money-pulse/services/users-service/internal/repository"
	"github.com/adil-faiyaz98/money-pulse/services/users-service/internal/service"
	"github.com/adil-faiyaz98/money-pulse/services/users-service/internal/transport/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger := logger.NewLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", err)
	}

	// Initialize token manager
	tokenManager, err := auth.NewTokenManager(cfg.JWTSecret)
	if err != nil {
		logger.Fatal("Failed to initialize token manager", err)
	}

	// Initialize repository
	userRepo, err := repository.NewUserRepository(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to initialize repository", err)
	}

	// Initialize service
	userService := service.NewUserService(userRepo, tokenManager)

	// Initialize HTTP handler
	userHandler := http.NewUserHandler(userService, tokenManager)

	// Initialize Gin router
	router := gin.Default()

	// Register routes
	userHandler.RegisterRoutes(router)

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", err)
	}

	logger.Info("Server exiting")
}
