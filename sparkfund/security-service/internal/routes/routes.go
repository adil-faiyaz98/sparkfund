package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"sparkfund/security-service/internal/handlers"
	"sparkfund/security-service/internal/metrics"
	"sparkfund/security-service/internal/middleware"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(router *gin.Engine, handler *handlers.Handler, metrics *metrics.Metrics) {
	// Middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.MetricsMiddleware(metrics))
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"time":    time.Now().Format(time.RFC3339),
			"service": "security-service",
		})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Security routes
		security := api.Group("/security")
		{
			security.POST("/validate-token", handler.ValidateToken)
			security.POST("/generate-token", handler.GenerateToken)
			security.POST("/refresh-token", handler.RefreshToken)
			// Add more security endpoints here
		}
	}
}
