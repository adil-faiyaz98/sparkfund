package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/investment-service/internal/config"
	"github.com/sparkfund/investment-service/internal/database"
	"github.com/sparkfund/investment-service/internal/handlers"
	"github.com/sparkfund/investment-service/internal/logger"
	"github.com/sparkfund/investment-service/internal/metrics"
	"github.com/sparkfund/investment-service/internal/repository"
	"github.com/sparkfund/investment-service/internal/routes"
	"github.com/sparkfund/investment-service/internal/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.joho/godotenv"
)

// @title Investment Management API
// @version 1.0
// @description API for managing client investments.

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Initialize logger
	log, err := logger.NewLogger(os.Getenv("ENV"))
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Connect to database
	db, err := database.NewDB(cfg.Database, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.CloseDB(db)

	// Initialize metrics
	metrics := metrics.NewAutoRegisterMetrics()

	// Initialize repository
	repo := repository.NewPostgresRepository(db, log)

	// Initialize service
	service := services.NewService(log, cfg, repo)

	// Initialize handler
	handler := handlers.NewHandler(log, service)

	// Initialize router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, handler)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := cfg.Port
	log.Info("Starting server", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
