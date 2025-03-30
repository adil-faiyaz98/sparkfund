package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/kyc-service/internal/config"
	"github.com/sparkfund/kyc-service/internal/handlers"
	"github.com/sparkfund/kyc-service/internal/middleware"
	"github.com/sparkfund/kyc-service/internal/security"
	"github.com/sparkfund/kyc-service/internal/server"
)

var (
	authConfig       security.AuthConfig
	validationConfig security.ValidationConfig
	encryptionConfig security.EncryptionConfig
	auditLogger      *security.AuditLogger
	handler          *handlers.Handler
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize security components
	authConfig = security.DefaultAuthConfig()
	validationConfig = security.DefaultValidationConfig()
	encryptionConfig = security.DefaultEncryptionConfig()
	auditConfig := security.DefaultAuditConfig()

	// Create audit logger
	auditLogger, err = security.NewAuditLogger(auditConfig)
	if err != nil {
		log.Fatalf("Failed to create audit logger: %v", err)
	}

	// Create handler instance
	handler = handlers.NewHandler(authConfig, validationConfig, encryptionConfig, auditLogger)

	// Create Gin router
	router := gin.New()

	// Add middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware(middleware.DefaultLogConfig()))
	router.Use(middleware.PrometheusMiddleware())
	router.Use(middleware.MetricsMiddleware())
	router.Use(middleware.HealthCheckMiddleware())
	router.Use(security.AuthMiddleware(authConfig))
	router.Use(security.AuditMiddleware(auditLogger))
	router.Use(middleware.SecurityMiddleware(middleware.DefaultSecurityConfig()))
	router.Use(middleware.RateLimitMiddleware(middleware.DefaultRateLimitConfig()))
	router.Use(middleware.CacheMiddleware(middleware.DefaultCacheConfig()))
	router.Use(middleware.GzipMiddleware())
	router.Use(middleware.TimeoutMiddleware(30 * time.Second))
	router.Use(middleware.ValidateRequestSize(10 * 1024 * 1024)) // 10MB
	router.Use(middleware.ValidateContentType([]string{"application/json", "multipart/form-data"}))

	// Add routes with role-based access control
	api := router.Group("/api/v1")
	{
		// Public routes
		api.POST("/auth/login", handler.HandleLogin)
		api.POST("/auth/refresh", handler.HandleRefreshToken)
		api.POST("/auth/register", handler.HandleRegister)

		// Protected routes
		kyc := api.Group("/kyc")
		kyc.Use(security.RequireRole("user", "admin"))
		{
			kyc.POST("", handler.HandleCreateKYC)
			kyc.GET("/:userID", handler.HandleGetKYC)
			kyc.PUT("", handler.HandleUpdateKYC)
			kyc.PATCH("/:userID/status", handler.HandleUpdateKYCStatus)
			kyc.GET("", handler.HandleListKYC)
			kyc.DELETE("/:userID", handler.HandleDeleteKYC)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(security.RequireRole("admin"))
		{
			admin.GET("/audit", handler.HandleGetAuditLogs)
			admin.GET("/metrics", handler.HandleGetMetrics)
			admin.POST("/config", handler.HandleUpdateConfig)
		}
	}

	// Health check and metrics endpoints
	router.GET("/health", handler.HandleHealth)
	router.GET("/metrics", handler.HandleMetrics)

	// Initialize server
	srv, err := server.New(cfg, router)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
