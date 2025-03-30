package main

import (
	"log"
	"os"
	"time"

	"sparkfund/services/kyc-service/internal/handler"
	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/repository"
	"sparkfund/services/kyc-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
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

	// Initialize router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(authMiddleware())

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

	// Health check
	router.GET("/health", healthCheck)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// TODO: Implement proper JWT validation
		// For now, just extract user ID from token
		userID := "123" // This should be extracted from the token
		c.Set("user_id", userID)

		c.Next()
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"time":   time.Now(),
	})
}
