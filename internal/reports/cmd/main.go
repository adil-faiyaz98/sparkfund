package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adil-faiyaz98/structgen/internal/reports/handlers"
	"github.com/adil-faiyaz98/structgen/internal/reports/repository/postgres"
	"github.com/adil-faiyaz98/structgen/internal/reports/service"
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
	repo := postgres.NewReportRepository(db)

	// Initialize service
	svc := service.NewReportService(repo)

	// Initialize handler
	handler := handlers.NewReportHandler(svc)

	// Initialize router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Define routes
	v1 := router.Group("/api/v1")
	{
		reports := v1.Group("/reports")
		{
			reports.POST("", handler.CreateReport)
			reports.GET("/:id", handler.GetReport)
			reports.GET("/user/:userId", handler.GetUserReports)
			reports.PUT("/:id/status", handler.UpdateReportStatus)
			reports.DELETE("/:id", handler.DeleteReport)
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
