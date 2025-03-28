package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sparkfund/api-gateway/internal/middleware"
	"github.com/sparkfund/api-gateway/internal/proxy"
)

func main() {
	// Set environment to development if not set
	if os.Getenv("ENV") == "" {
		os.Setenv("ENV", "development")
	}

	// Initialize security middleware
	securityConfig := middleware.SecurityConfig{
		RateLimit: struct {
			RequestsPerMinute int
			BurstSize         int
		}{
			RequestsPerMinute: 100,
			BurstSize:         200,
		},
		DoS: struct {
			MaxHeaderSize    int64
			MaxBodySize      int64
			MaxConnections   int
			ConnectionWindow time.Duration
		}{
			MaxHeaderSize:    1 << 20,  // 1MB
			MaxBodySize:      10 << 20, // 10MB
			MaxConnections:   100,
			ConnectionWindow: time.Minute,
		},
		JWT: struct {
			SecretKey     []byte
			TokenExpiry   time.Duration
			RefreshExpiry time.Duration
		}{
			SecretKey:     []byte(os.Getenv("JWT_SECRET")),
			TokenExpiry:   time.Hour * 24,
			RefreshExpiry: time.Hour * 24 * 7,
		},
	}

	securityMiddleware := middleware.NewSecurityMiddleware(securityConfig)

	// Set up Gin router
	router := gin.Default()

	// Set trusted proxies
	router.SetTrustedProxies([]string{
		"127.0.0.1",
		"172.16.0.0/12",
		"172.17.0.0/12",
		"172.18.0.0/12",
		"172.19.0.0/16",
		"192.168.0.0/16",
		"10.0.0.0/8",
	})

	// Add after setting trusted proxies
	router.ForwardedByClientIP = true
	router.RemoteIPHeaders = []string{"X-Forwarded-For", "X-Real-IP"}

	// First register metrics endpoint (before security middleware)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API info endpoint
	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "SparkFund API Gateway",
			"version": "1.0",
			"status":  "running",
			"endpoints": []string{
				"/api/v1/investments",
				"/api/v1/portfolios",
				"/api/v1/transactions",
			},
		})
	})

	// Service discovery endpoint
	router.GET("/api/v1/investment-service", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "Investment Service",
			"status":  "running",
			"url":     "http://investment-service:8081",
		})
	})

	// Health check (without auth)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Add Prometheus server to IP whitelist
	securityMiddleware.AddToIPWhitelist("172.18.0.4") // Prometheus IP
	securityMiddleware.AddToIPWhitelist("172.18.0.1") // Local testing IP

	// Apply security middleware
	securityMiddleware.Apply(router)

	// Investment service routes
	investments := router.Group("/api/v1/investments")
	{
		investments.POST("/", proxy.ProxyToInvestmentService)
		investments.GET("/:id", proxy.ProxyToInvestmentService)
		investments.GET("/", proxy.ProxyToInvestmentService)
		investments.PUT("/:id", proxy.ProxyToInvestmentService)
		investments.DELETE("/:id", proxy.ProxyToInvestmentService)
	}

	// Portfolio routes
	portfolios := router.Group("/api/v1/portfolios")
	{
		portfolios.POST("/", proxy.ProxyToInvestmentService)
		portfolios.GET("/:id", proxy.ProxyToInvestmentService)
		portfolios.PUT("/:id", proxy.ProxyToInvestmentService)
		portfolios.DELETE("/:id", proxy.ProxyToInvestmentService)
	}

	// Transaction routes
	transactions := router.Group("/api/v1/transactions")
	{
		transactions.POST("/", proxy.ProxyToInvestmentService)
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("API Gateway starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
