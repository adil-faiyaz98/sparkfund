package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/api-gateway/internal/middleware"
	"github.com/sparkfund/api-gateway/internal/proxy"
)

func main() {
	// Set up Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Create security config with proper fields matching your SecurityConfig struct
	securityConfig := middleware.SecurityConfig{
		// Initialize with default values for required fields
		RateLimit: struct {
			RequestsPerMinute int
			BurstSize         int
		}{
			RequestsPerMinute: 60,
			BurstSize:         10,
		},
		DoS: struct {
			MaxHeaderSize    int64
			MaxBodySize      int64
			MaxConnections   int
			ConnectionWindow time.Duration
		}{
			MaxHeaderSize:    8192,
			MaxBodySize:      10 * 1024 * 1024, // 10MB
			MaxConnections:   100,
			ConnectionWindow: 60 * time.Second,
		},
		JWT: struct {
			SecretKey     []byte
			TokenExpiry   time.Duration
			RefreshExpiry time.Duration
		}{
			SecretKey:     []byte("your-secret-key"),
			TokenExpiry:   24 * time.Hour,
			RefreshExpiry: 7 * 24 * time.Hour,
		},
	}

	securityMiddleware := middleware.NewSecurityMiddleware(securityConfig)
	securityMiddleware.AddToIPWhitelist("172.18.0.4")

	// Add Prometheus IP to whitelist
	securityMiddleware.AddToIPWhitelist("172.18.0.4") // Prometheus server IP

	// Add a public metrics endpoint (before applying security middleware)
	router.GET("/metrics", func(c *gin.Context) {
		c.String(http.StatusOK, "# HELP api_requests_total Total API requests\n# TYPE api_requests_total counter\napi_requests_total 100\n")
	})

	securityMiddleware.Apply(router)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

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
