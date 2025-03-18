package api

import (
	"github.com/gin-gonic/gin"

	"github.com/adilm/money-pulse/internal/pkg/middleware"
)

// SetupRouter configures the API routes
func SetupRouter(
	handler *Handler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(authMiddleware.AuthRequired())
	{
		// User routes
		api.GET("/me", handler.GetCurrentUser)
		api.PUT("/me", handler.UpdateCurrentUser)

		// Transaction routes
		transactions := api.Group("/transactions")
		{
			transactions.GET("", handler.GetTransactions)
			transactions.POST("", handler.CreateTransaction)
			transactions.GET("/:id", handler.GetTransaction)
			transactions.PUT("/:id", handler.UpdateTransaction)
			transactions.DELETE("/:id", handler.DeleteTransaction)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", handler.GetCategories)
			categories.POST("", handler.CreateCategory)
			categories.PUT("/:id", handler.UpdateCategory)
			categories.DELETE("/:id", handler.DeleteCategory)
		}

		// Account routes
		accounts := api.Group("/accounts")
		{
			accounts.GET("", handler.GetAccounts)
			accounts.POST("", handler.CreateAccount)
			accounts.PUT("/:id", handler.UpdateAccount)
			accounts.DELETE("/:id", handler.DeleteAccount)
		}

		// Dashboard routes
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/summary", handler.GetDashboardSummary)
		}
	}

	return router
}
