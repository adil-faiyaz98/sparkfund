package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/investment-service/internal/models"
	"github.com/sparkfund/investment-service/internal/services"
	"go.uber.org/zap"
)

type Handler struct {
	logger  *zap.Logger
	service services.Service
}

func NewHandler(logger *zap.Logger, service services.Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

// GetInvestmentRecommendation godoc
// @Summary Get AI-based investment recommendation
// @Description Get AI-based investment recommendation for a specific investment.
// @ID getInvestmentRecommendation
// @Produce json
// @Param investmentId path string true "Investment ID"
// @Success 200 {object} models.InvestmentRecommendation
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /investments/{investmentId}/recommendation [get]
func (h *Handler) GetInvestmentRecommendation(c *gin.Context) {
	investmentId := c.Param("investmentId")
	h.logger.Info("Getting investment recommendation", zap.String("investmentId", investmentId))

	recommendation, err := h.service.GetInvestmentRecommendation(c.Request.Context(), investmentId)
	if err != nil {
		h.logger.Error("Failed to get investment recommendation", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to get investment recommendation"})
		return
	}

	c.JSON(http.StatusOK, recommendation)
}

// GetClientInvestments godoc
// @Summary Get client investments
// @Description Retrieves all investments for a client.
// @ID getClientInvestments
// @Produce json
// @Param clientId path string true "Client ID"
// @Success 200 {object} models.Investment
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /clients/{clientId}/investments [get]
func (h *Handler) GetClientInvestments(c *gin.Context) {
	clientId := c.Param("clientId")
	h.logger.Info("Getting investments for client", zap.String("clientId", clientId))

	investments, err := h.service.GetClientInvestments(c.Request.Context(), clientId)
	if err != nil {
		h.logger.Error("Failed to get investments", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to get investments"})
		return
	}

	c.JSON(http.StatusOK, investments)
}

// CreatePortfolio godoc
// @Summary Create a portfolio
// @Description Creates a new portfolio for a client.
// @ID createPortfolio
// @Accept json
// @Produce json
// @Param portfolio body models.Portfolio true "Portfolio details"
// @Success 201 {object} models.Portfolio
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /portfolios [post]
func (h *Handler) CreatePortfolio(c *gin.Context) {
	var portfolio models.Portfolio
	if err := c.BindJSON(&portfolio); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{Code: 400, Message: "Invalid request body"})
		return
	}

	createdPortfolio, err := h.service.CreatePortfolio(c.Request.Context(), portfolio)
	if err != nil {
		h.logger.Error("Failed to create portfolio", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to create portfolio"})
		return
	}

	c.JSON(http.StatusCreated, createdPortfolio)
}

// GetPortfolio godoc
// @Summary Get a portfolio by ID
// @Description Retrieves details for a specific portfolio.
// @ID getPortfolioById
// @Produce json
// @Param portfolioId path string true "Portfolio ID"
// @Success 200 {object} models.Portfolio
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /portfolios/{portfolioId} [get]
func (h *Handler) GetPortfolio(c *gin.Context) {
	portfolioId := c.Param("portfolioId")
	h.logger.Info("Getting portfolio", zap.String("portfolioId", portfolioId))

	portfolio, err := h.service.GetPortfolio(c.Request.Context(), portfolioId)
	if err != nil {
		h.logger.Error("Failed to get portfolio", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to get portfolio"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// UpdatePortfolio godoc
// @Summary Update a portfolio
// @Description Updates details of a specific portfolio.
// @ID updatePortfolio
// @Accept json
// @Produce json
// @Param portfolioId path string true "Portfolio ID"
// @Param portfolio body models.Portfolio true "Portfolio object that needs to be updated"
// @Success 200 {object} models.Portfolio
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /portfolios/{portfolioId} [put]
func (h *Handler) UpdatePortfolio(c *gin.Context) {
	portfolioId := c.Param("portfolioId")
	var portfolio models.Portfolio
	if err := c.BindJSON(&portfolio); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{Code: 400, Message: "Invalid request body"})
		return
	}

	updatedPortfolio, err := h.service.UpdatePortfolio(c.Request.Context(), portfolioId, portfolio)
	if err != nil {
		h.logger.Error("Failed to update portfolio", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to update portfolio"})
		return
	}

	c.JSON(http.StatusOK, updatedPortfolio)
}

// DeletePortfolio godoc
// @Summary Delete a portfolio
// @Description Deletes a specific portfolio.
// @ID deletePortfolio
// @Param portfolioId path string true "Portfolio ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /portfolios/{portfolioId} [delete]
func (h *Handler) DeletePortfolio(c *gin.Context) {
	portfolioId := c.Param("portfolioId")
	h.logger.Info("Deleting portfolio", zap.String("portfolioId", portfolioId))

	err := h.service.DeletePortfolio(c.Request.Context(), portfolioId)
	if err != nil {
		h.logger.Error("Failed to delete portfolio", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to delete portfolio"})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateInvestment godoc
// @Summary Create a new investment
// @Description Creates a new investment. Requires client ID and portfolio ID.
// @ID createInvestment
// @Accept json
// @Produce json
// @Param investment body models.Investment true "Investment details"
// @Success 201 {object} models.Investment
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /investments [post]
func (h *Handler) CreateInvestment(c *gin.Context) {
	var investment models.Investment
	if err := c.BindJSON(&investment); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{Code: 400, Message: "Invalid request body"})
		return
	}

	createdInvestment, err := h.service.CreateInvestment(c.Request.Context(), investment)
	if err != nil {
		h.logger.Error("Failed to create investment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to create investment"})
		return
	}

	c.JSON(http.StatusCreated, createdInvestment)
}

// GetInvestment godoc
// @Summary Get an investment by ID
// @Description Retrieves details for a specific investment.
// @ID getInvestmentById
// @Produce json
// @Param investmentId path string true "Investment ID"
// @Success 200 {object} models.Investment
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /investments/{investmentId} [get]
func (h *Handler) GetInvestment(c *gin.Context) {
	investmentId := c.Param("investmentId")
	h.logger.Info("Getting investment", zap.String("investmentId", investmentId))

	investment, err := h.service.GetInvestment(c.Request.Context(), investmentId)
	if err != nil {
		h.logger.Error("Failed to get investment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to get investment"})
		return
	}

	c.JSON(http.StatusOK, investment)
}

// UpdateInvestment godoc
// @Summary Update an investment
// @Description Updates details of a specific investment.
// @ID updateInvestment
// @Accept json
// @Produce json
// @Param investmentId path string true "Investment ID"
// @Param investment body models.Investment true "Investment object that needs to be updated"
// @Success 200 {object} models.Investment
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /investments/{investmentId} [put]
func (h *Handler) UpdateInvestment(c *gin.Context) {
	investmentId := c.Param("investmentId")
	var investment models.Investment
	if err := c.BindJSON(&investment); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{Code: 400, Message: "Invalid request body"})
		return
	}

	updatedInvestment, err := h.service.UpdateInvestment(c.Request.Context(), investmentId, investment)
	if err != nil {
		h.logger.Error("Failed to update investment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to update investment"})
		return
	}

	c.JSON(http.StatusOK, updatedInvestment)
}

// DeleteInvestment godoc
// @Summary Delete an investment
// @Description Liquidates a specific investment.
// @ID deleteInvestment
// @Param investmentId path string true "Investment ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /investments/{investmentId} [delete]
func (h *Handler) DeleteInvestment(c *gin.Context) {
	investmentId := c.Param("investmentId")
	h.logger.Info("Deleting investment", zap.String("investmentId", investmentId))

	err := h.service.DeleteInvestment(c.Request.Context(), investmentId)
	if err != nil {
		h.logger.Error("Failed to delete investment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to delete investment"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetInvestmentForecast godoc
// @Summary Get AI-based forecast for an investment
// @Description Get AI-based forecast for a specific investment.
// @ID getInvestmentForecast
// @Produce json
// @Param investmentId path string true "Investment ID"
// @Success 200 {object} models.InvestmentForecast
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /investments/{investmentId}/forecast [get]
func (h *Handler) GetInvestmentForecast(c *gin.Context) {
	investmentId := c.Param("investmentId")
	h.logger.Info("Getting forecast for investment", zap.String("investmentId", investmentId))

	forecast, err := h.service.GetInvestmentForecast(c.Request.Context(), investmentId)
	if err != nil {
		h.logger.Error("Failed to get investment forecast", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to get investment forecast"})
		return
	}

	c.JSON(http.StatusOK, forecast)
}

// GetClientWill godoc
// @Summary Get client will
// @Description Retrieve the will details for a specific client.
// @ID getClientWill
// @Produce json
// @Param clientId path string true "Client ID"
// @Success 200 {object} models.Will
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /clients/{clientId}/wills [get]
func (h *Handler) GetClientWill(c *gin.Context) {
	clientId := c.Param("clientId")
	h.logger.Info("Getting will for client", zap.String("clientId", clientId))

	will, err := h.service.GetClientWill(c.Request.Context(), clientId)
	if err != nil {
		h.logger.Error("Failed to get will", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to get will"})
		return
	}

	c.JSON(http.StatusOK, will)
}

// CreateOrUpdateClientWill godoc
// @Summary Create or update a client will
// @Description Creates or updates a will specifying the distribution of a client's investments.
// @ID createOrUpdateClientWill
// @Accept json
// @Produce json
// @Param clientId path string true "Client ID"
// @Param will body models.Will true "Will details to create or update"
// @Success 200 {object} models.Will
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /clients/{clientId}/wills [post]
func (h *Handler) CreateOrUpdateClientWill(c *gin.Context) {
	clientId := c.Param("clientId")
	var will models.Will
	if err := c.BindJSON(&will); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{Code: 400, Message: "Invalid request body"})
		return
	}

	updatedWill, err := h.service.CreateOrUpdateClientWill(c.Request.Context(), clientId, will)
	if err != nil {
		h.logger.Error("Failed to create or update will", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to create or update will"})
		return
	}

	c.JSON(http.StatusOK, updatedWill)
}

// GetClientsWithdrawalThreshold godoc
// @Summary Get client withdrawal threshold
// @Description Retrieve the automated withdrawal threshold for a client.
// @ID getClientsWithdrawalThreshold
// @Produce json
// @Param clientId path string true "Client ID"
// @Success 200 {object} models.WithdrawalThreshold
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /clients/{clientId}/withdrawalThreshold [get]
func (h *Handler) GetClientsWithdrawalThreshold(c *gin.Context) {
	clientId := c.Param("clientId")
	h.logger.Info("Getting withdrawal threshold for client", zap.String("clientId", clientId))

	threshold, err := h.service.GetClientsWithdrawalThreshold(c.Request.Context(), clientId)
	if err != nil {
		h.logger.Error("Failed to get withdrawal threshold", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to get withdrawal threshold"})
		return
	}

	c.JSON(http.StatusOK, threshold)
}

// SetClientWithdrawalThreshold godoc
// @Summary Set client withdrawal threshold
// @Description Sets the automated withdrawal threshold for a specific client.
// @ID setClientWithdrawalThreshold
// @Accept json
// @Produce json
// @Param clientId path string true "Client ID"
// @Param threshold body models.WithdrawalThreshold true "The automated withdrawal threshold to set"
// @Success 200 {object} models.WithdrawalThreshold
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /clients/{clientId}/withdrawalThreshold [post]
func (h *Handler) SetClientWithdrawalThreshold(c *gin.Context) {
	clientId := c.Param("clientId")
	var threshold models.WithdrawalThreshold
	if err := c.BindJSON(&threshold); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{Code: 400, Message: "Invalid request body"})
		return
	}

	updatedThreshold, err := h.service.SetClientWithdrawalThreshold(c.Request.Context(), clientId, threshold)
	if err != nil {
		h.logger.Error("Failed to set withdrawal threshold", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to set withdrawal threshold"})
		return
	}

	c.JSON(http.StatusOK, updatedThreshold)
}

// UpdateClientWithdrawalThreshold godoc
// @Summary Update client withdrawal threshold
// @Description Updates the automated withdrawal threshold for a client.
// @ID updateClientWithdrawalThreshold
// @Accept json
// @Produce json
// @Param clientId path string true "Client ID"
// @Param threshold body models.WithdrawalThreshold true "The amount at which to trigger an automated withdrawal"
// @Success 200 {object} models.WithdrawalThreshold
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /clients/{clientId}/withdrawalThreshold [put]
func (h *Handler) UpdateClientWithdrawalThreshold(c *gin.Context) {
	clientId := c.Param("clientId")
	var threshold models.WithdrawalThreshold
	if err := c.BindJSON(&threshold); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.Error{Code: 400, Message: "Invalid request body"})
		return
	}

	updatedThreshold, err := h.service.UpdateClientWithdrawalThreshold(c.Request.Context(), clientId, threshold)
	if err != nil {
		h.logger.Error("Failed to update withdrawal threshold", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.Error{Code: 500, Message: "Failed to update withdrawal threshold"})
		return
	}

	c.JSON(http.StatusOK, updatedThreshold)
}
