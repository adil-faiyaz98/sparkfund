package routes

import (
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, handler *handlers.Handler) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Email routes
		emails := v1.Group("/emails")
		{
			emails.POST("", handler.SendEmail)
			emails.GET("", handler.GetEmailLogs)
		}

		// Template routes
		templates := v1.Group("/templates")
		{
			templates.POST("", handler.CreateTemplate)
			templates.GET("/:id", handler.GetTemplate)
			templates.PUT("/:id", handler.UpdateTemplate)
			templates.DELETE("/:id", handler.DeleteTemplate)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
