package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sparkfund/kyc-service/internal/config"
)

// Verification represents a verification record
type Verification struct {
	ID          uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID      uuid.UUID  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001" gorm:"type:uuid;not null"`
	KYCID       uuid.UUID  `json:"kyc_id" example:"550e8400-e29b-41d4-a716-446655440002" gorm:"type:uuid;not null"`
	DocumentID  *uuid.UUID `json:"document_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440003" gorm:"type:uuid"`
	VerifierID  *uuid.UUID `json:"verifier_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440004" gorm:"type:uuid"`
	Method      string     `json:"method" example:"DOCUMENT" gorm:"type:varchar(50);not null"`
	Status      string     `json:"status" example:"PENDING" gorm:"type:varchar(50);not null"`
	Notes       string     `json:"notes,omitempty" example:"Verification in progress" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at" example:"2025-04-05T15:04:05Z" gorm:"not null;default:now()"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2025-04-05T15:04:05Z" gorm:"not null;default:now()"`
	CompletedAt *time.Time `json:"completed_at,omitempty" example:"2025-04-05T15:04:05Z"`
}

// VerificationRequest represents a request to create or update a verification
type VerificationRequest struct {
	UserID     uuid.UUID  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001" binding:"required"`
	KYCID      uuid.UUID  `json:"kyc_id" example:"550e8400-e29b-41d4-a716-446655440002" binding:"required"`
	DocumentID *uuid.UUID `json:"document_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440003"`
	Method     string     `json:"method" example:"DOCUMENT" binding:"required"`
	Status     string     `json:"status" example:"PENDING" binding:"required"`
	Notes      string     `json:"notes,omitempty" example:"Verification in progress"`
}

// VerificationResponse represents a response containing a verification
type VerificationResponse struct {
	Verification Verification `json:"verification"`
}

// VerificationsResponse represents a response containing multiple verifications
type VerificationsResponse struct {
	Verifications []Verification `json:"verifications"`
	Total         int            `json:"total" example:"10"`
	Page          int            `json:"page" example:"1"`
	PageSize      int            `json:"page_size" example:"20"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad request"`
	Details string `json:"details,omitempty" example:"Invalid verification ID format"`
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Setup router
	router := gin.Default()

	// Register routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": cfg.App.Name,
			"version": cfg.App.Version,
			"env":     cfg.App.Environment,
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Verification routes
		verifications := api.Group("/verifications")
		{
			verifications.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "List verifications endpoint",
				})
			})

			verifications.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Get verification endpoint",
					"id":      c.Param("id"),
				})
			})

			verifications.POST("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Create verification endpoint",
				})
			})

			verifications.PUT("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Update verification endpoint",
					"id":      c.Param("id"),
				})
			})

			verifications.DELETE("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Delete verification endpoint",
					"id":      c.Param("id"),
				})
			})
		}
	}

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
