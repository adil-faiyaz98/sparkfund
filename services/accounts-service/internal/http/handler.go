package http

import (
	"net/http"
	"strconv"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/service"
	"github.com/gin-gonic/gin"
)

// AccountHandler handles HTTP requests for account operations
type AccountHandler struct {
	service *service.AccountService
}

// NewAccountHandler creates a new HTTP handler for accounts
func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

// RegisterRoutes registers all route handlers
func (h *AccountHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		accounts := api.Group("/accounts")
		{
			accounts.POST("", h.createAccount)
			accounts.GET("/:id", h.getAccount)
			accounts.PUT("/:id", h.updateAccountStatus)
			accounts.DELETE("/:id", h.deleteAccount)
			accounts.GET("/user/:userID", h.listUserAccounts)
			accounts.POST("/:id/deposit", h.depositFunds)
			accounts.POST("/:id/withdraw", h.withdrawFunds)
			accounts.POST("/transfer", h.transferFunds)
		}
	}
}

// Request/response structures
type createAccountRequest struct {
	UserID         string  `json:"user_id" binding:"required"`
	AccountType    string  `json:"account_type" binding:"required"`
	Currency       string  `json:"currency" binding:"required"`
	InitialDeposit float64 `json:"initial_deposit"`
}

type updateAccountStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type fundTransferRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

type transferFundsRequest struct {
	SourceID      string  `json:"source_id" binding:"required"`
	DestinationID string  `json:"destination_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Description   string  `json:"description"`
}

// Handler methods
func (h *AccountHandler) createAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.service.CreateAccount(c, req.UserID, req.AccountType, req.Currency, req.InitialDeposit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) getAccount(c *gin.Context) {
	id := c.Param("id")
	account, err := h.service.GetAccount(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) updateAccountStatus(c *gin.Context) {
	id := c.Param("id")
	var req updateAccountStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.service.UpdateAccountStatus(c, id, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) deleteAccount(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteAccount(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AccountHandler) listUserAccounts(c *gin.Context) {
	userID := c.Param("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	accounts, total, err := h.service.ListUserAccounts(c, userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *AccountHandler) depositFunds(c *gin.Context) {
	id := c.Param("id")
	var req fundTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.service.DepositFunds(c, id, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) withdrawFunds(c *gin.Context) {
	id := c.Param("id")
	var req fundTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.service.WithdrawFunds(c, id, req.Amount)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrInsufficientFunds {
			statusCode = http.StatusBadRequest
		} else if err == service.ErrAccountNotActive {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) transferFunds(c *gin.Context) {
	var req transferFundsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.service.TransferFunds(c, req.SourceID, req.DestinationID, req.Amount)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrInsufficientFunds {
			statusCode = http.StatusBadRequest
		} else if err == service.ErrAccountNotActive {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}
