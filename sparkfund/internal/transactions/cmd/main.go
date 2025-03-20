package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adil-faiyaz98/sparkfund/internal/transactions/handlers"
	"github.com/adil-faiyaz98/sparkfund/internal/transactions/repository/postgres"
	"github.com/adil-faiyaz98/sparkfund/internal/transactions/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize database connection
	db, err := sqlx.Connect("postgres", "postgres://postgres:postgres@localhost:5432/transactions?sslmode=disable")
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repository
	repo := postgres.NewRepository(db)

	// Initialize service
	svc := service.NewTransactionService(repo)

	// Initialize handler
	handler := handlers.NewHandler(svc)

	// Initialize router
	router := gin.New()

	// Setup routes
	handler.RegisterRoutes(router)

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
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
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
