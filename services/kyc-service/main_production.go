package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adil-faiyaz98/sparkfund/pkg/config"
	"github.com/adil-faiyaz98/sparkfund/pkg/logger"
	"github.com/adil-faiyaz98/sparkfund/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	log := logger.GetLogger()
	defer log.Sync()

	// Load configuration
	if err := config.Load("./config"); err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}
	cfg := config.Get()

	// Set Gin mode
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	
	// Add security middleware
	corsConfig := middleware.DefaultCORSConfig()
	corsConfig.AllowedOrigins = cfg.Security.AllowedOrigins
	corsConfig.AllowedMethods = cfg.Security.AllowedMethods
	corsConfig.AllowedHeaders = cfg.Security.AllowedHeaders
	router.Use(middleware.CORS(corsConfig))
	router.Use(middleware.SecurityHeaders())
	
	// Add rate limiting
	rateLimitConfig := middleware.DefaultRateLimiterConfig()
	rateLimitConfig.Enabled = cfg.RateLimit.Enabled
	rateLimitConfig.Requests = cfg.RateLimit.Requests
	rateLimitConfig.Window = cfg.RateLimit.Window
	rateLimitConfig.Burst = cfg.RateLimit.Burst
	router.Use(middleware.RateLimiter(rateLimitConfig))
	
	// Add JWT authentication if enabled
	jwtConfig := middleware.DefaultJWTConfig()
	jwtConfig.Secret = cfg.JWT.Secret
	jwtConfig.Enabled = cfg.JWT.Enabled
	router.Use(middleware.JWTAuth(jwtConfig))
	
	// Add CSRF protection if enabled
	if cfg.Security.EnableCSRF {
		router.Use(middleware.CSRFProtection())
	}
	
	// Add circuit breaker
	cbConfig := middleware.DefaultCircuitBreakerConfig()
	cbConfig.Enabled = cfg.CircuitBreaker.Enabled
	cbConfig.Timeout = cfg.CircuitBreaker.Timeout
	cbConfig.MaxConcurrentReqs = cfg.CircuitBreaker.MaxConcurrentReqs
	cbConfig.ErrorThresholdPerc = cfg.CircuitBreaker.ErrorThresholdPerc
	cbConfig.RequestVolumeThresh = uint64(cfg.CircuitBreaker.RequestVolumeThresh)
	cbConfig.SleepWindow = cfg.CircuitBreaker.SleepWindow
	router.Use(middleware.CircuitBreakerMiddleware(cbConfig))
	
	// Add tracing
	tracingConfig := middleware.DefaultTracingConfig()
	tracingConfig.ServiceName = cfg.Tracing.ServiceName
	tracingConfig.Enabled = cfg.Tracing.Enabled
	router.Use(middleware.TracingMiddleware(tracingConfig))
	
	// Add metrics endpoint
	if cfg.Metrics.Enabled {
		router.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}
	
	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "kyc-service",
			"version": "1.0.0",
		})
	})
	
	// Add API endpoints
	v1 := router.Group("/api/v1")
	{
		// Authentication endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/login", loginHandler)
		}
		
		// KYC verification endpoints
		verifications := v1.Group("/verifications")
		{
			verifications.POST("", createVerificationHandler)
			verifications.GET("", listVerificationsHandler)
			verifications.GET("/:id", getVerificationHandler)
			verifications.PUT("/:id", updateVerificationHandler)
		}
		
		// AI integration endpoints
		ai := v1.Group("/ai")
		{
			ai.GET("/models", getAIModelsHandler)
			ai.POST("/analyze-document", analyzeDocumentHandler)
			ai.POST("/match-faces", matchFacesHandler)
			ai.POST("/analyze-risk", analyzeRiskHandler)
			ai.POST("/detect-anomalies", detectAnomaliesHandler)
			ai.POST("/process-document", processDocumentHandler)
		}
		
		// API key endpoint
		v1.GET("/get-api-key", getAPIKeyHandler)
	}
	
	// Serve Swagger UI
	router.GET("/swagger-ui.html", func(c *gin.Context) {
		c.File("swagger-ui.html")
	})
	
	// Serve Swagger JSON
	router.GET("/swagger.json", func(c *gin.Context) {
		c.File("swagger.json")
	})
	
	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	
	// Start server in a goroutine
	go func() {
		logger.Info("Starting KYC service", 
			logger.String("port", cfg.Server.Port),
			logger.String("environment", cfg.Environment))
			
		var err error
		if cfg.TLS.Enabled {
			err = server.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		} else {
			err = server.ListenAndServe()
		}
		
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", logger.ErrorField(err))
		}
	}()
	
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	// Shutdown server gracefully
	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", logger.ErrorField(err))
	}
	
	logger.Info("Server exited gracefully")
}

// Handler implementations
func loginHandler(c *gin.Context) {
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"token": "sample-jwt-token",
		"user": gin.H{
			"id":         "123e4567-e89b-12d3-a456-426614174000",
			"email":      "user@example.com",
			"first_name": "John",
			"last_name":  "Doe",
			"role":       "user",
		},
	})
}

func createVerificationHandler(c *gin.Context) {
	// Implementation details
	c.JSON(http.StatusCreated, gin.H{
		"verification": gin.H{
			"id":            "123e4567-e89b-12d3-a456-426614174001",
			"user_id":       "123e4567-e89b-12d3-a456-426614174000",
			"kyc_id":        "123e4567-e89b-12d3-a456-426614174002",
			"document_id":   "123e4567-e89b-12d3-a456-426614174003",
			"method":        "AI",
			"status":        "PENDING",
			"created_at":    time.Now().Format(time.RFC3339),
			"updated_at":    time.Now().Format(time.RFC3339),
		},
	})
}

func listVerificationsHandler(c *gin.Context) {
	// Get pagination parameters
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"verifications": []gin.H{
			{
				"id":            "123e4567-e89b-12d3-a456-426614174001",
				"user_id":       "123e4567-e89b-12d3-a456-426614174000",
				"kyc_id":        "123e4567-e89b-12d3-a456-426614174002",
				"document_id":   "123e4567-e89b-12d3-a456-426614174003",
				"method":        "AI",
				"status":        "PENDING",
				"created_at":    time.Now().Format(time.RFC3339),
				"updated_at":    time.Now().Format(time.RFC3339),
			},
		},
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      1,
			"total_pages": 1,
		},
	})
}

func getVerificationHandler(c *gin.Context) {
	// Get verification ID
	id := c.Param("id")
	
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"verification": gin.H{
			"id":            id,
			"user_id":       "123e4567-e89b-12d3-a456-426614174000",
			"kyc_id":        "123e4567-e89b-12d3-a456-426614174002",
			"document_id":   "123e4567-e89b-12d3-a456-426614174003",
			"method":        "AI",
			"status":        "PENDING",
			"created_at":    time.Now().Format(time.RFC3339),
			"updated_at":    time.Now().Format(time.RFC3339),
		},
	})
}

func updateVerificationHandler(c *gin.Context) {
	// Get verification ID
	id := c.Param("id")
	
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"verification": gin.H{
			"id":            id,
			"user_id":       "123e4567-e89b-12d3-a456-426614174000",
			"kyc_id":        "123e4567-e89b-12d3-a456-426614174002",
			"document_id":   "123e4567-e89b-12d3-a456-426614174003",
			"method":        "AI",
			"status":        "APPROVED",
			"created_at":    time.Now().Format(time.RFC3339),
			"updated_at":    time.Now().Format(time.RFC3339),
		},
	})
}

func getAIModelsHandler(c *gin.Context) {
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"models": []gin.H{
			{
				"id":              "123e4567-e89b-12d3-a456-426614174010",
				"name":            "Document Verification Model",
				"version":         "1.0.0",
				"type":            "DOCUMENT",
				"accuracy":        0.98,
				"last_trained_at": time.Now().Format(time.RFC3339),
				"created_at":      time.Now().Format(time.RFC3339),
				"updated_at":      time.Now().Format(time.RFC3339),
			},
			{
				"id":              "123e4567-e89b-12d3-a456-426614174011",
				"name":            "Face Recognition Model",
				"version":         "1.0.0",
				"type":            "FACE",
				"accuracy":        0.95,
				"last_trained_at": time.Now().Format(time.RFC3339),
				"created_at":      time.Now().Format(time.RFC3339),
				"updated_at":      time.Now().Format(time.RFC3339),
			},
		},
	})
}

func analyzeDocumentHandler(c *gin.Context) {
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"id":            "123e4567-e89b-12d3-a456-426614174020",
		"verification_id": "123e4567-e89b-12d3-a456-426614174001",
		"document_id":   "123e4567-e89b-12d3-a456-426614174003",
		"document_type": "PASSPORT",
		"is_authentic":  true,
		"confidence":    0.95,
		"extracted_data": gin.H{
			"full_name":       "John Smith",
			"document_number": "X123456789",
			"date_of_birth":   "1990-01-01",
			"expiry_date":     "2030-01-01",
			"issuing_country": "United States",
		},
		"issues":      []string{},
		"created_at":  time.Now().Format(time.RFC3339),
	})
}

func matchFacesHandler(c *gin.Context) {
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"id":              "123e4567-e89b-12d3-a456-426614174021",
		"verification_id": "123e4567-e89b-12d3-a456-426614174001",
		"document_id":     "123e4567-e89b-12d3-a456-426614174003",
		"selfie_id":       "123e4567-e89b-12d3-a456-426614174004",
		"is_match":        true,
		"confidence":      0.92,
		"created_at":      time.Now().Format(time.RFC3339),
	})
}

func analyzeRiskHandler(c *gin.Context) {
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"id":              "123e4567-e89b-12d3-a456-426614174022",
		"verification_id": "123e4567-e89b-12d3-a456-426614174001",
		"user_id":         "123e4567-e89b-12d3-a456-426614174000",
		"risk_score":      0.15,
		"risk_level":      "LOW",
		"risk_factors":    []string{},
		"device_info": gin.H{
			"ip_address":    "192.168.1.1",
			"user_agent":    "Mozilla/5.0",
			"device_type":   "Desktop",
			"os":            "Windows",
			"browser":       "Chrome",
			"location":      "New York, USA",
			"captured_time": time.Now().Format(time.RFC3339),
		},
		"ip_address":    "192.168.1.1",
		"location":      "New York, USA",
		"created_at":    time.Now().Format(time.RFC3339),
	})
}

func detectAnomaliesHandler(c *gin.Context) {
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"id":              "123e4567-e89b-12d3-a456-426614174023",
		"verification_id": "123e4567-e89b-12d3-a456-426614174001",
		"user_id":         "123e4567-e89b-12d3-a456-426614174000",
		"is_anomaly":      false,
		"anomaly_score":   0.05,
		"anomaly_type":    nil,
		"reasons":         []string{},
		"device_info": gin.H{
			"ip_address":    "192.168.1.1",
			"user_agent":    "Mozilla/5.0",
			"device_type":   "Desktop",
			"os":            "Windows",
			"browser":       "Chrome",
			"location":      "New York, USA",
			"captured_time": time.Now().Format(time.RFC3339),
		},
		"created_at":    time.Now().Format(time.RFC3339),
	})
}

func processDocumentHandler(c *gin.Context) {
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"id":           "123e4567-e89b-12d3-a456-426614174001",
		"status":       "COMPLETED",
		"notes":        "All verification checks passed",
		"completed_at": time.Now().Format(time.RFC3339),
	})
}

func getAPIKeyHandler(c *gin.Context) {
	// Implementation details
	c.JSON(http.StatusOK, gin.H{
		"api_key": "your-api-key",
		"note":    "Use this API key in the X-API-Key header when calling the AI service",
	})
}
