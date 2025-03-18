package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adil-faiyaz98/structgen/internal/transactions/handlers"
	"github.com/adil-faiyaz98/structgen/internal/transactions/repository/postgres"
	"github.com/adil-faiyaz98/structgen/internal/transactions/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	// Initialize database connection
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	repo := postgres.NewTransactionRepository(db)

	// Initialize service
	svc := service.NewTransactionService(repo)

	// Initialize handler
	handler := handlers.NewTransactionHandler(svc)

	// Initialize router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Define routes
	v1 := router.Group("/api/v1")
	{
		transactions := v1.Group("/transactions")
		{
			transactions.POST("", handler.CreateTransaction)
			transactions.GET("/:id", handler.GetTransaction)
			transactions.GET("/user/:userId", handler.GetUserTransactions)
			transactions.GET("/account/:accountId", handler.GetAccountTransactions)
			transactions.PUT("/:id/status", handler.UpdateTransactionStatus)
			transactions.DELETE("/:id", handler.DeleteTransaction)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown server
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
