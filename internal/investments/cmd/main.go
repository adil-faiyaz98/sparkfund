package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"your-project/internal/investments/handlers"
	"your-project/internal/investments/repository"
	"your-project/internal/investments/service"

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
	repo := repository.NewPostgresRepository(db)

	// Initialize service
	svc := service.NewInvestmentService(repo)

	// Initialize handler
	handler := handlers.NewInvestmentHandler(svc)

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		investments := api.Group("/investments")
		{
			investments.POST("", handler.CreateInvestment)
			investments.GET("/:id", handler.GetInvestment)
			investments.GET("/user/:userId", handler.GetUserInvestments)
			investments.GET("/account/:accountId", handler.GetAccountInvestments)
			investments.PUT("/:id", handler.UpdateInvestment)
			investments.DELETE("/:id", handler.DeleteInvestment)
			investments.GET("/symbol/:symbol", handler.GetInvestmentsBySymbol)
			investments.PUT("/:id/price/:price", handler.UpdateInvestmentPrice)
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
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
