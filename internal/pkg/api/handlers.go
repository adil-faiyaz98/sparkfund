package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/adilm/money-pulse/internal/pkg/database"
	"github.com/adilm/money-pulse/internal/pkg/models"
)

// Handler contains all HTTP handlers
type Handler struct {
	transactionRepo *database.TransactionRepository
	userRepo        *database.UserRepository
	categoryRepo    *database.CategoryRepository
	accountRepo     *database.AccountRepository
}

// NewHandler creates a new Handler
func NewHandler(
	transactionRepo *database.TransactionRepository,
	userRepo *database.UserRepository,
	categoryRepo *database.CategoryRepository,
	accountRepo *database.AccountRepository,
) *Handler {
	return &Handler{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		categoryRepo:    categoryRepo,
		accountRepo:     accountRepo,
	}
}

// CreateTransaction handles transaction creation
func (h *Handler) CreateTransaction(c *gin.Context) {
	var req struct {
		Amount      float64 `json:"amount" binding:"required"`
		Description string  `json:"description"`
		Date        string  `json:"date" binding:"required"`
		CategoryID  uint    `json:"category_id"`
		AccountID   uint    `json:"account_id" binding:"required"`
		Type        string  `json:"type" binding:"required,oneof=income expense"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := getUserIDFromContext(c)

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	transaction := models.Transaction{
		UserID:      userID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        date,
		CategoryID:  req.CategoryID,
		AccountID:   req.AccountID,
		Type:        req.Type,
	}

	if err := h.transactionRepo.CreateTransaction(c.Request.Context(), &transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetTransactions handles retrieving transactions
func (h *Handler) GetTransactions(c *gin.Context) {
	userID := getUserIDFromContext(c)

	filter := &database.TransactionFilter{}

	if categoryID := c.Query("category_id"); categoryID != "" {
		id, err := strconv.ParseUint(categoryID, 10, 32)
		if err == nil {
			filter.CategoryID = uint(id)
		}
	}

	if accountID := c.Query("account_id"); accountID != "" {
		id, err := strconv.ParseUint(accountID, 10, 32)
		if err == nil {
			filter.AccountID = uint(id)
		}
	}

	if startDate := c.Query("start_date"); startDate != "" {
		date, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			filter.StartDate = date
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		date, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			filter.EndDate = date
		}
	}

	if transactionType := c.Query("type"); transactionType != "" {
		filter.Type = transactionType
	}

	transactions, err := h.transactionRepo.GetTransactions(c.Request.Context(), userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// getUserIDFromContext extracts userID from gin context
func getUserIDFromContext(c *gin.Context) uint {
	id, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	userID, ok := id.(uint)
	if !ok {
		return 0
	}

	return userID
}
