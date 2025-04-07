package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/adil-faiyaz98/sparkfund/services/service-template/internal/api/handlers"
	"github.com/adil-faiyaz98/sparkfund/services/service-template/internal/api/middleware"
	"github.com/adil-faiyaz98/sparkfund/services/service-template/internal/database"
	"github.com/adil-faiyaz98/sparkfund/services/service-template/internal/service"
)

// SetupRoutes sets up all the routes for the API
func SetupRoutes(router *gin.RouterGroup, db *database.Database, logger *logrus.Logger) {
	// Create services
	exampleService := service.NewExampleService(db, logger)

	// Create handlers
	exampleHandler := handlers.NewExampleHandler(exampleService, logger)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(logger)
	loggerMiddleware := middleware.NewLoggerMiddleware(logger)

	// Apply global middleware
	router.Use(loggerMiddleware.LogRequest())

	// Example routes
	examples := router.Group("/examples")
	{
		examples.GET("", exampleHandler.GetAll)
		examples.GET("/:id", exampleHandler.GetByID)
		examples.POST("", authMiddleware.Authenticate(), exampleHandler.Create)
		examples.PUT("/:id", authMiddleware.Authenticate(), exampleHandler.Update)
		examples.DELETE("/:id", authMiddleware.Authenticate(), exampleHandler.Delete)
	}
}
