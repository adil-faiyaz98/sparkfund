package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sparkfund/security-monitoring/internal/ai"
	"github.com/sparkfund/security-monitoring/internal/config"
	"github.com/sparkfund/security-monitoring/internal/handlers"
	"github.com/sparkfund/security-monitoring/internal/middleware"
	"github.com/sparkfund/security-monitoring/internal/monitoring"
	"github.com/sparkfund/security-monitoring/internal/security"
)

func main() {
	// Initialize logger
	logger := setupLogger()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize AI components
	aiEngine := ai.NewEngine(cfg.AI)

	// Initialize security monitoring components
	securityMonitor := security.NewMonitor(cfg.Security, aiEngine)
	threatDetector := security.NewThreatDetector(cfg.Security, aiEngine)
	intrusionDetector := security.NewIntrusionDetector(cfg.Security, aiEngine)
	malwareDetector := security.NewMalwareDetector(cfg.Security, aiEngine)
	patternAnalyzer := security.NewPatternAnalyzer(cfg.Security, aiEngine)

	// Initialize monitoring
	metrics := monitoring.NewMetrics()
	alertManager := monitoring.NewAlertManager(cfg.Alerts)

	// Set up Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.RateLimit(cfg.RateLimit))

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Security monitoring endpoints
	api := router.Group("/api/v1")
	{
		security := api.Group("/security")
		{
			security.POST("/analyze", handlers.AnalyzeSecurity(securityMonitor))
			security.GET("/threats", handlers.GetThreats(threatDetector))
			security.GET("/incidents", handlers.GetIncidents(intrusionDetector))
			security.POST("/scan", handlers.ScanMalware(malwareDetector))
			security.GET("/patterns", handlers.GetPatterns(patternAnalyzer))
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")
}

func setupLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}
