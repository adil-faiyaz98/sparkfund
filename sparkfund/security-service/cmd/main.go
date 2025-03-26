package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"sparkfund/security-service/internal/config"
	"sparkfund/security-service/internal/database"
	"sparkfund/security-service/internal/handlers"
	"sparkfund/security-service/internal/logger"
	"sparkfund/security-service/internal/metrics"
	"sparkfund/security-service/internal/repositories"
	"sparkfund/security-service/internal/routes"
	"sparkfund/security-service/internal/services"
)

// @title Security Service API
// @version 1.0
// @description API for security-related operations.

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
	metricsClient := metrics.NewAutoRegisterMetrics()

	// Initialize repository
	repo := repositories.NewRepository(db, log)

	// Initialize service
	service := services.NewService(log, cfg, repo)

	// Initialize handler
	handler := handlers.NewHandler(log, service)

	// Initialize router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, handler, metricsClient)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := cfg.Port
	log.Info("Starting security service", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
