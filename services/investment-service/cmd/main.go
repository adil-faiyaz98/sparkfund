package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "investment-service/docs"
	"investment-service/internal/config"
	"investment-service/internal/database"
	"investment-service/internal/handlers"
	"investment-service/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Investment Service API
// @version         1.0
// @description     A service for managing investments and portfolios.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
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
	router.Use(middleware.SecurityHeaders())
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

	// Set up routes
	investments := api.Group("/investments")
	{
		investments.POST("/", handlers.CreateInvestment)
		investments.GET("/:id", handlers.GetInvestment)
		investments.GET("/", handlers.ListInvestments)
		investments.PUT("/:id", handlers.UpdateInvestment)
		investments.DELETE("/:id", handlers.DeleteInvestment)
	}

	portfolios := api.Group("/portfolios")
	{
		portfolios.POST("/", handlers.CreatePortfolio)
		portfolios.GET("/:id", handlers.GetPortfolio)
		portfolios.PUT("/:id", handlers.UpdatePortfolio)
		portfolios.DELETE("/:id", handlers.DeletePortfolio)
	}

	transactions := api.Group("/transactions")
	{
		transactions.POST("/", handlers.CreateTransaction)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Infof("Starting server on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited gracefully")
}
