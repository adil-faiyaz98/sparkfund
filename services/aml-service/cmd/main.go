package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aml-service/internal/config"
	"aml-service/internal/database"
	"aml-service/internal/handlers"
	"aml-service/internal/metrics"
	"aml-service/internal/middleware"
	"aml-service/internal/repositories"
	"aml-service/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize metrics
	metrics := metrics.NewMetrics()

	// Initialize database connection
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	amlRepo := repositories.NewAMLRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Initialize services
	amlService := services.NewAMLService(amlRepo, transactionRepo, logger)

	// Initialize router
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.Metrics(metrics))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Security.CORS.AllowedOrigins,
		AllowedMethods:   cfg.Security.CORS.AllowedMethods,
		AllowedHeaders:   cfg.Security.CORS.AllowedHeaders,
		ExposedHeaders:   cfg.Security.CORS.ExposedHeaders,
		AllowCredentials: cfg.Security.CORS.AllowCredentials,
		MaxAge:           cfg.Security.CORS.MaxAge,
	}))

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))
	})

	// Metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/aml", handlers.NewAMLHandler(amlService, logger).Routes())
	})

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server
	go func() {
		logger.Info("Starting server", zap.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
