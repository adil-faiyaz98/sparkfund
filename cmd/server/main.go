package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sparkfund/services/kyc-service/internal/handler"
	"github.com/sparkfund/services/kyc-service/internal/logger"
	"github.com/sparkfund/services/kyc-service/internal/metrics"
	"github.com/sparkfund/services/kyc-service/internal/middleware"
	"github.com/sparkfund/services/kyc-service/internal/model"
	"github.com/sparkfund/services/kyc-service/internal/repository"
	"github.com/sparkfund/services/kyc-service/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize logger
	if err := logger.InitLogger(os.Getenv("LOG_LEVEL")); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.GetLogger().Sync()

	// Initialize database
	db, err := initDB()
	if err != nil {
		logger.Fatal("Failed to initialize database", logger.ErrorField(err))
	}

	// Initialize repositories
	documentRepo := repository.NewDocumentRepository(db)
	verificationRepo := repository.NewVerificationRepository(db)

	// Initialize services
	documentService := service.NewDocumentService(documentRepo, verificationRepo, os.Getenv("UPLOAD_DIR"))
	verificationService := service.NewVerificationService(verificationRepo, documentRepo)

	// Initialize handlers
	documentHandler := handler.NewDocumentHandler(documentService)
	verificationHandler := handler.NewVerificationHandler(verificationService)

	// Initialize rate limiter
	rateLimiter, err := middleware.NewRateLimiter(os.Getenv("REDIS_URL"))
	if err != nil {
		logger.Fatal("Failed to initialize rate limiter", logger.ErrorField(err))
	}

	// Initialize router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(middleware.AuthMiddleware())
	router.Use(rateLimiter.RateLimit(100, time.Minute)) // 100 requests per minute
	router.Use(metrics.PrometheusMiddleware())

	// API routes
	api := router.Group("/api/v1")
	{
		// Document routes
		documents := api.Group("/documents")
		{
			documents.POST("", documentHandler.UploadDocument)
			documents.GET("/:id", documentHandler.GetDocument)
			documents.GET("/user/:user_id", documentHandler.GetUserDocuments)
			documents.PUT("/:id/status", documentHandler.UpdateDocumentStatus)
			documents.DELETE("/:id", documentHandler.DeleteDocument)
			documents.GET("/:id/history", documentHandler.GetDocumentHistory)
			documents.GET("/stats", documentHandler.GetDocumentStats)
			documents.GET("/:id/summary", documentHandler.GetDocumentSummary)
			documents.GET("/expired", documentHandler.GetExpiredDocuments)
			documents.GET("/pending", documentHandler.GetPendingDocuments)
			documents.GET("/rejected", documentHandler.GetRejectedDocuments)
			documents.GET("/type/:type", documentHandler.GetDocumentsByType)
			documents.GET("/date-range", documentHandler.GetDocumentsByDateRange)
			documents.PUT("/:id/metadata", documentHandler.UpdateDocumentMetadata)
		}

		// Verification routes
		verifications := api.Group("/verifications")
		{
			verifications.POST("", verificationHandler.CreateVerification)
			verifications.GET("/:id", verificationHandler.GetVerification)
			verifications.GET("/document/:document_id", verificationHandler.GetDocumentVerifications)
			verifications.PUT("/:id/status", verificationHandler.UpdateVerificationStatus)
			verifications.DELETE("/:id", verificationHandler.DeleteVerification)
			verifications.GET("/:id/history", verificationHandler.GetVerificationHistory)
			verifications.GET("/stats", verificationHandler.GetVerificationStats)
			verifications.GET("/document/:document_id/summary", verificationHandler.GetVerificationSummary)
			verifications.GET("/expired", verificationHandler.GetExpiredVerifications)
			verifications.GET("/pending", verificationHandler.GetPendingVerifications)
			verifications.GET("/failed", verificationHandler.GetFailedVerifications)
			verifications.GET("/verifier/:verifier_id", verificationHandler.GetVerificationsByVerifier)
			verifications.GET("/method/:method", verificationHandler.GetVerificationsByMethod)
			verifications.GET("/date-range", verificationHandler.GetVerificationsByDateRange)
			verifications.PUT("/:id/confidence", verificationHandler.UpdateConfidenceScore)
			verifications.PUT("/:id/metadata", verificationHandler.UpdateVerificationMetadata)
		}
	}

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check
	router.GET("/health", healthCheck)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting server", logger.String("port", port))
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", logger.ErrorField(err))
	}
}

func initDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=sparkfund port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	err = db.AutoMigrate(
		&model.Document{},
		&model.DocumentHistory{},
		&model.Verification{},
		&model.VerificationHistory{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"time":   time.Now(),
	})
}
