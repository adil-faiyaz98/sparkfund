package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adil-faiyaz98/structgen/internal/loans/handlers"
	"github.com/adil-faiyaz98/structgen/internal/loans/repository/postgres"
	"github.com/adil-faiyaz98/structgen/internal/loans/service"
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
	repo := postgres.NewLoanRepository(db)

	// Initialize service
	svc := service.NewLoanService(repo)

	// Initialize handler
	handler := handlers.NewLoanHandler(svc)

	// Initialize router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Define routes
	v1 := router.Group("/api/v1")
	{
		loans := v1.Group("/loans")
		{
			loans.POST("", handler.CreateLoan)
			loans.GET("/:id", handler.GetLoan)
			loans.GET("/user/:userId", handler.GetUserLoans)
			loans.GET("/account/:accountId", handler.GetAccountLoans)
			loans.PUT("/:id/status", handler.UpdateLoanStatus)
			loans.POST("/:id/payments", handler.MakePayment)
			loans.GET("/:id/payments", handler.GetLoanPayments)
			loans.DELETE("/:id", handler.DeleteLoan)
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
