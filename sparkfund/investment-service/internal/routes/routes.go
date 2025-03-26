package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sparkfund/investment-service/internal/handlers"
)

func SetupRoutes(router *gin.Engine, handler *handlers.Handler) {
	v1 := router.Group("/api/v1")
	{
		clients := v1.Group("/clients/:clientId")
		{
			clients.GET("/investments", handler.GetClientInvestments)
		}

		portfolios := v1.Group("/portfolios")
		{
			portfolios.POST("", handler.CreatePortfolio)
			portfolios.GET("/:portfolioId", handler.GetPortfolio)
			portfolios.PUT("/:portfolioId", handler.UpdatePortfolio)
			portfolios.DELETE("/:portfolioId", handler.DeletePortfolio)
		}

		investments := v1.Group("/investments")
		{
			investments.GET("/:investmentId", handler.GetInvestment)
			investments.PUT("/:investmentId", handler.UpdateInvestment)
			investments.DELETE("/:investmentId", handler.DeleteInvestment)
			investments.GET("/:investmentId/recommendation", handler.GetInvestmentRecommendation) // New route
		}
	}
}
