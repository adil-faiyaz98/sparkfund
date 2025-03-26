package routes

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sparkfund/email-service/internal/handlers"
	"github.com/sparkfund/email-service/internal/middleware"
	"github.com/sparkfund/email-service/internal/services"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, handler *handlers.Handler, authService *services.AuthService, logger *zap.Logger) {
	// Middleware
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestID())
	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.ErrorHandler(logger))

	// CORS configuration
	config := cors.DefaultConfig()

	// Allow all origins during development, otherwise, use a comma-separated list from env
	if os.Getenv("APP_ENV") == "development" {
		config.AllowAllOrigins = true
	} else {
		allowedOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
		if allowedOrigins != "" {
			config.AllowOrigins = strings.Split(allowedOrigins, ",")
		} else {
			logger.Warn("CORS_ALLOW_ORIGINS not set, using default (empty) list.  Requests may be blocked.")
		}
	}

	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowCredentials = true // If you need to allow cookies, set this to true
	router.Use(cors.New(config))

	// Rate limiter
	limit := getEnvIntOrDefault("RATE_LIMIT", 100)       // Requests per second
	burst := getEnvIntOrDefault("RATE_LIMIT_BURST", 200) // Max burst size
	router.Use(middleware.RateLimiter(rate.Limit(limit), burst))

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Email routes
		emails := v1.Group("/emails")
		emails.Use(middleware.Auth(authService)) // Require authentication
		{
			emails.POST("", handler.SendEmail)
			emails.GET("", handler.GetEmailLogs)
		}

		// Template routes
		templates := v1.Group("/templates")
		templates.Use(middleware.Auth(authService)) // Require authentication
		{
			templates.POST("", handler.CreateTemplate)
			templates.GET("/:id", handler.GetTemplate)
			templates.PUT("/:id", handler.UpdateTemplate)
			templates.DELETE("/:id", handler.DeleteTemplate)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue // Or log the error
	}
	return intValue
}

// use export CORS_ALLOW_ORIGINS="https://your-frontend.com,https://another-frontend.com"
