package main

import (
	"log"
	"os"

	_ "investment-service/docs" // This is where the generated docs will be
	"investment-service/internal/database"
	"investment-service/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Investment Service API
// @version         1.0
// @description     A service for managing investments and portfolios.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /api/v1
// @schemes   http

func main() {
	// Initialize database
	database.InitDB()

	// Initialize Prometheus metrics
	buildInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_build_info",
			Help: "Information about the Go build.",
		},
		[]string{"version", "revision", "branch", "goversion"},
	)

	// Register Prometheus collectors
	prometheus.MustRegister(buildInfo)

	// Set up Gin router with proper logging
	router := gin.New()
	router.Use(gin.Logger()) // Add logger middleware to see requests
	router.Use(gin.Recovery())

	// Set trusted proxies explicitly
	router.SetTrustedProxies([]string{
		"127.0.0.1",
		"172.16.0.0/12",
		"172.17.0.0/12",
		"172.18.0.0/12",
		"172.19.0.0/16",
		"192.168.0.0/16",
		"10.0.0.0/8",
	})

	// Public endpoints
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// API routes (no authentication for local development)
	api := router.Group("/api/v1")
	{
		// Investment routes
		investments := api.Group("/investments")
		{
			investments.POST("/", handlers.CreateInvestment)
			investments.GET("/:id", handlers.GetInvestment)
			investments.GET("/", handlers.ListInvestments)
			investments.PUT("/:id", handlers.UpdateInvestment)
			investments.DELETE("/:id", handlers.DeleteInvestment)
		}

		// Portfolio routes
		portfolios := api.Group("/portfolios")
		{
			portfolios.POST("/", handlers.CreatePortfolio)
			portfolios.GET("/:id", handlers.GetPortfolio)
			portfolios.PUT("/:id", handlers.UpdatePortfolio)
			portfolios.DELETE("/:id", handlers.DeletePortfolio)
		}

		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.POST("/", handlers.CreateTransaction)
		}
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Make sure this matches the port in docker-compose
	}

	// Start server on the correct port
	log.Printf("Investment Service starting on port %s", port)
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
