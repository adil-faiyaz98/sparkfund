package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/investment-service/internal/database"
	"github.com/sparkfund/investment-service/internal/handlers"
)

func main() {
	// Initialize database
	database.InitDB()

	// Set up Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Investment routes
	investments := router.Group("/api/v1/investments")
	{
		investments.POST("/", handlers.CreateInvestment)
		investments.GET("/:id", handlers.GetInvestment)
		investments.GET("/", handlers.ListInvestments)
		investments.PUT("/:id", handlers.UpdateInvestment)
		investments.DELETE("/:id", handlers.DeleteInvestment)
	}

	// Portfolio routes
	portfolios := router.Group("/api/v1/portfolios")
	{
		portfolios.POST("/", handlers.CreatePortfolio)
		portfolios.GET("/:id", handlers.GetPortfolio)
		portfolios.PUT("/:id", handlers.UpdatePortfolio)
		portfolios.DELETE("/:id", handlers.DeletePortfolio)
	}

	// Transaction routes
	transactions := router.Group("/api/v1/transactions")
	{
		transactions.POST("/", handlers.CreateTransaction)
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 