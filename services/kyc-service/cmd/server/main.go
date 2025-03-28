package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/kyc-service/config"
	"github.com/sparkfund/kyc-service/internal/handler"
	"github.com/sparkfund/kyc-service/internal/repository"
	"github.com/sparkfund/kyc-service/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	kycRepo := repository.NewKYCRepository(db)

	// Initialize services
	kycService := service.NewKYCService(kycRepo)

	// Initialize handlers
	kycHandler := handler.NewKYCHandler(kycService)

	// Initialize router
	router := gin.Default()

	// Initialize routes
	initializeRoutes(router, kycHandler)

	// Create server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
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
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func initializeRoutes(router *gin.Engine, kycHandler *handler.KYCHandler) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Register KYC routes
	kycHandler.RegisterRoutes(router)
}

func handleKYCVerification(c *gin.Context) {
	// TODO: Implement KYC verification logic
	c.JSON(http.StatusOK, gin.H{
		"message": "KYC verification endpoint",
	})
}

func handleKYCStatus(c *gin.Context) {
	// TODO: Implement KYC status check logic
	c.JSON(http.StatusOK, gin.H{
		"message": "KYC status endpoint",
	})
}
