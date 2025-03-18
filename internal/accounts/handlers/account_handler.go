package handlers

import (
	"net/http"

	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/accounts"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountHandler struct {
	service accounts.AccountService
}

func NewAccountHandler(service accounts.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

// CreateAccount handles the creation of a new account
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var account accounts.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateAccount(c.Request.Context(), &account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccount retrieves an account by ID
func (h *AccountHandler) GetAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	account, err := h.service.GetAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetUserAccounts retrieves all accounts for a user
func (h *AccountHandler) GetUserAccounts(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	accounts, err := h.service.GetUserAccounts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// UpdateAccount updates an existing account
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var account accounts.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account.ID = id
	if err := h.service.UpdateAccount(c.Request.Context(), &account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// DeleteAccount deletes an account
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	if err := h.service.DeleteAccount(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAccountByNumber retrieves an account by account number
func (h *AccountHandler) GetAccountByNumber(c *gin.Context) {
	accountNumber := c.Param("accountNumber")
	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account number is required"})
		return
	}

	account, err := h.service.GetAccountByNumber(c.Request.Context(), accountNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}
