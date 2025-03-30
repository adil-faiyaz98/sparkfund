package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "kyc-service/docs"
	"kyc-service/internal/config"
	"kyc-service/internal/database"
	"kyc-service/internal/handlers"
	"kyc-service/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           KYC Service API
// @version         1.0
// @description     A service for managing KYC verifications and reviews.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8082
// @BasePath  /api/v1
// @schemes   http https

var version = "development" // Replaced during build

func main() {
	// Set up logging
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	cfg := config.Get()

	// Set log level from config
	logLevel, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Warnf("Invalid log level %s, defaulting to info", cfg.Log.Level)
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set up router with middlewares
	router := gin.New()

	// Add middlewares in correct order
	router.Use(middleware.RequestLogger())
	router.Use(gin.Recovery())
	router.Use(middleware.SecurityHeaders(middleware.DefaultSecurityConfig()))
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter())

	// Health endpoints (no auth required)
	router.GET("/health", handlers.HealthCheck)
	router.GET("/live", handlers.LivenessCheck)
	router.GET("/ready", handlers.ReadinessCheck)

	// Metrics endpoint
	if cfg.Metrics.Enabled {
		router.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}

	// Swagger docs (disable in production)
	if os.Getenv("APP_ENV") != "production" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API routes (with auth)
	api := router.Group("/api/v1")
	api.Use(middleware.JWTAuth())

	// KYC routes
	kyc := api.Group("/kyc")
	{
		kyc.POST("/verifications", handlers.CreateKYCVerification)
		kyc.GET("/verifications/:id", handlers.GetKYCVerification)
		kyc.PUT("/verifications/:id", handlers.UpdateKYCVerification)
		kyc.GET("/verifications", handlers.ListKYCVerifications)
		kyc.POST("/verifications/:id/review", handlers.ReviewKYCVerification)
		kyc.GET("/documents/:id", handlers.GetKYCDocument)
		kyc.POST("/documents", handlers.UploadKYCDocument)
		kyc.GET("/watchlist", handlers.ListKYCWatchlist)
		kyc.POST("/watchlist", handlers.AddToWatchlist)
		kyc.DELETE("/watchlist/:id", handlers.RemoveFromWatchlist)
	}

	// Set up server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exiting")
}
