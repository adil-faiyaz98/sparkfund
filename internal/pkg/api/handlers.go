package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/adil-faiyaz98/structgen/internal/pkg/database"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Check if account exists and belongs to user
	account, err := h.accountRepo.GetByID(req.AccountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if account.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account doesn't belong to user"})
		return
	}

	// Create transaction
	transaction := &models.Transaction{
		Amount:      req.Amount,
		Description: req.Description,
		Date:        date,
		CategoryID:  req.CategoryID,
		AccountID:   req.AccountID,
		Type:        req.Type,
	}

	if err := h.transactionRepo.Create(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetTransactions handles retrieving transactions
func (h *Handler) GetTransactions(c *gin.Context) {
	userID := getUserIDFromContext(c)
	
	// Parse query parameters
	accountID := c.Query("account_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	// Apply filters
	filters := map[string]interface{}{}
	if accountID != "" {
		accountIDInt, err := strconv.ParseUint(accountID, 10, 64)
		if err == nil {
			filters["account_id"] = uint(accountIDInt)
			
			// Verify account belongs to user
			account, err := h.accountRepo.GetByID(uint(accountIDInt))
			if err != nil || account.UserID != userID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this account"})
				return
			}
		}
	} else {
		// If no account specified, get all user accounts
		accounts, err := h.accountRepo.GetByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accounts"})
			return
		}
		
		accountIDs := make([]uint, len(accounts))
		for i, acc := range accounts {
			accountIDs[i] = acc.ID
		}
		
		filters["account_ids"] = accountIDs
	}
	
	// Date filtering
	if startDate != "" {
		parsedStartDate, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			filters["start_date"] = parsedStartDate
		}
	}
	
	if endDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			filters["end_date"] = parsedEndDate
		}
	}
	
	// Get transactions with pagination
	transactions, total, err := h.transactionRepo.GetFiltered(filters, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetTransaction handles retrieving a single transaction by ID
func (h *Handler) GetTransaction(c *gin.Context) {
	userID := getUserIDFromContext(c)
	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}
	
	transaction, err := h.transactionRepo.GetByID(uint(transactionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	
	// Verify transaction's account belongs to user
	account, err := h.accountRepo.GetByID(transaction.AccountID)
	if err != nil || account.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	
	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction handles updating a transaction
func (h *Handler) UpdateTransaction(c *gin.Context) {
	userID := getUserIDFromContext(c)
	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}
	
	// Check if transaction exists
	transaction, err := h.transactionRepo.GetByID(uint(transactionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	
	// Verify transaction's account belongs to user
	account, err := h.accountRepo.GetByID(transaction.AccountID)
	if err != nil || account.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	
	// Parse request body
	var req struct {
		Amount      *float64 `json:"amount"`
		Description *string  `json:"description"`
		Date        *string  `json:"date"`
		CategoryID  *uint    `json:"category_id"`
		Type        *string  `json:"type"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Update fields if provided
	if req.Amount != nil {
		transaction.Amount = *req.Amount
	}
	
	if req.Description != nil {
		transaction.Description = *req.Description
	}
	
	if req.Date != nil {
		date, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		transaction.Date = date
	}
	
	if req.CategoryID != nil {
		transaction.CategoryID = *req.CategoryID
	}
	
	if req.Type != nil {
		if *req.Type != "income" && *req.Type != "expense" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be 'income' or 'expense'"})
			return
		}
		transaction.Type = *req.Type
	}
	
	// Save changes
	if err := h.transactionRepo.Update(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}
	
	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction handles deleting a transaction
func (h *Handler) DeleteTransaction(c *gin.Context) {
	userID := getUserIDFromContext(c)
	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}
	
	// Check if transaction exists
	transaction, err := h.transactionRepo.GetByID(uint(transactionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	
	// Verify transaction's account belongs to user
	account, err := h.accountRepo.GetByID(transaction.AccountID)
	if err != nil || account.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	
	// Delete transaction
	if err := h.transactionRepo.Delete(uint(transactionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

// getUserIDFromContext extracts userID from gin context
func getUserIDFromContext(c *gin.Context) uint {
	id, exists := c.Get("user_id")
	if !exists {
		return 0
	}
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.userRepo.Create(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	} Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	// Don't return the password hash
	user.Password = ""`json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
	c.JSON(http.StatusCreated, user)ame" binding:"required"`
}	LastName  string `json:"last_name" binding:"required"`
	}
// Login handles user authentication
func (h *Handler) Login(c *gin.Context) { nil {
	var req struct {tusBadRequest, gin.H{"error": err.Error()})
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`

	// Check if user already exists
	existingUser, err := h.userRepo.GetByEmail(c.Request.Context(), req.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	// Create new user
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process registration"})
		return
	}

	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.userRepo.Create(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	// Don't return the password hash
	user.Password = ""

	c.JSON(http.StatusCreated, user)
}

// Login handles user authentication
func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	user, err := h.userRepo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Verify password
	if !checkPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Don't return the password hash
	user.Password = ""
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return": token,
	}"user":  user,
	})
	// Get existing user
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil || user == nil {eving the current authenticated user
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return:= getUserIDFromContext(c)
	}
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	// Update fields if provided
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Don't return the password hash
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// UpdateCurrentUser handles updating the current user's profile
func (h *Handler) UpdateCurrentUser(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing user
	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" && req.Email != user.Email {
		// Check if email is already taken
		existingUser, err := h.userRepo.GetByEmail(c.Request.Context(), req.Email)
		if err == nil && existingUser != nil && existingUser.ID != user.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
			return
		}
		user.Email = req.Email
	}

	if err := h.userRepo.Update(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// Don't return the password hash
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// GetTransaction handles retrieving a single transaction by ID
func (h *Handler) GetTransaction(c *gin.Context) {
	userID := getUserIDFromContext(c)
	transactionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}
	
	transaction, err := h.transactionRepo.GetByID(uint(transactionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	
	// Verify transaction's account belongs to user
	account, err := h.accountRepo.GetByID(transaction.AccountID)
	if err != nil || account.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	
    if err != nil {
	categories, err := h.categoryRepo.GetCategories(c.Request.Context(), userID): "Category not found"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories"})    }
		return
	}
ss denied"})
	c.JSON(http.StatusOK, categories)
}
  
func (h *Handler) CreateCategory(c *gin.Context) {    var req struct {
	userID := getUserIDFromContext(c)

	var req struct {
		Name  string `json:"name" binding:"required"`  
		Color string `json:"color"`    if err := c.ShouldBindJSON(&req); err != nil {
		Icon  string `json:"icon"`adRequest, gin.H{"error": err.Error()})
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}req.Name
  }
	category := models.Category{    
		Name:   req.Name,
		Color:  req.Color,
		Icon:   req.Icon,
		UserID: userID,  
	}    // Save changes
category); err != nil {
	if err := h.categoryRepo.CreateCategory(c.Request.Context(), &category); err != nil {       c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})        return



















































}	c.JSON(http.StatusCreated, account)	}		return		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})	if err := h.accountRepo.CreateAccount(c.Request.Context(), &account); err != nil {	}		UserID:         userID,		Currency:       req.Currency,		CurrentBalance: req.InitialBalance,		InitialBalance: req.InitialBalance,		Type:           req.Type,		Name:           req.Name,	account := models.Account{	}		return		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})	if err := c.ShouldBindJSON(&req); err != nil {	}		Currency       string  `json:"currency" binding:"required"`		InitialBalance float64 `json:"initial_balance"`		Type           string  `json:"type" binding:"required"`		Name           string  `json:"name" binding:"required"`	var req struct {	userID := getUserIDFromContext(c)func (h *Handler) CreateAccount(c *gin.Context) {}	c.JSON(http.StatusOK, accounts)	}		return		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch accounts"})	if err != nil {	accounts, err := h.accountRepo.GetAccounts(c.Request.Context(), userID)	userID := getUserIDFromContext(c)func (h *Handler) GetAccounts(c *gin.Context) {// Account handlers}	c.JSON(http.StatusCreated, category)	}		return

































































































































































































}    c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})        }        return        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})    if err := h.accountRepo.Delete(uint(accountID)); err != nil {    // Delete account        }        return        })            "action": "deactivate",            "error": "Cannot delete account with transactions. Consider deactivating it instead.",        c.JSON(http.StatusConflict, gin.H{    if hasTransactions {        }        return        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check account usage"})    if err != nil {    hasTransactions, err := h.transactionRepo.HasAccountTransactions(account.ID)    // Check if account has transactions        }        return        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})    if account.UserID != userID {        }        return        c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})    if err != nil {    account, err := h.accountRepo.GetByID(uint(accountID))    // Check if account exists and belongs to user        }        return        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})    if err != nil {    accountID, err := strconv.ParseUint(c.Param("id"), 10, 64)    userID := getUserIDFromContext(c)func (h *Handler) DeleteAccount(c *gin.Context) {// DeleteAccount handles deleting an account}    })        "updated_at": account.UpdatedAt,        "created_at": account.CreatedAt,        "is_active":  account.IsActive,        "currency":   account.Currency,        "balance":    balance,        "color":      account.Color,        "type":       account.Type,        "name":       account.Name,        "id":         account.ID,    c.JSON(http.StatusOK, gin.H{        }        balance = 0    if err != nil {    balance, err := h.transactionRepo.GetAccountBalance(account.ID)    // Get current balance        }        return        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})    if err := h.accountRepo.Update(account); err != nil {    // Save changes        }        account.IsActive = *req.IsActive    if req.IsActive != nil {        }        account.Color = *req.Color    if req.Color != nil {        }        account.Name = *req.Name    if req.Name != nil {    // Update fields if provided        }        return        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})    if err := c.ShouldBindJSON(&req); err != nil {        }        IsActive *bool   `json:"is_active"`        Color    *string `json:"color"`        Name     *string `json:"name"`    var req struct {        }        return        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})    if account.UserID != userID {        }        return        c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})    if err != nil {    account, err := h.accountRepo.GetByID(uint(accountID))    // Check if account exists and belongs to user



}	c.JSON(http.StatusCreated, account)	}		return		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})


        }        return        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
	if err := h.accountRepo.CreateAccount(c.Request.Context(), &account); err != nil {
    if err != nil {    accountID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	}		UserID:         userID,		Currency:       req.Currency,

    userID := getUserIDFromContext(c)func (h *Handler) UpdateAccount(c *gin.Context) {		CurrentBalance: req.InitialBalance,		InitialBalance: req.InitialBalance,

// UpdateAccount handles updating an account}    })        "created_at":     account.CreatedAt,


		Type:           req.Type,		Name:           req.Name,	account := models.Account{	}		return


        "is_active":      account.IsActive,        "currency":       account.Currency,        "balance":        req.InitialBalance,		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})	if err := c.ShouldBindJSON(&req); err != nil {	}		Currency       string  `json:"currency" binding:"required"`		InitialBalance float64 `json:"initial_balance"`		Type           string  `json:"type" binding:"required"`


        "color":          account.Color,        "type":           account.Type,


        "name":           account.Name,        "id":             account.ID,    c.JSON(http.StatusCreated, gin.H{
        }        return        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})

		Name           string  `json:"name" binding:"required"`	var req struct {	userID := getUserIDFromContext(c)

    if err != nil {    err := h.accountRepo.CreateWithInitialBalance(account, req.InitialBalance)func (h *Handler) CreateAccount(c *gin.Context) {}	c.JSON(http.StatusOK, accounts)




    // Create account in transaction        }        IsActive: true,        UserID:   userID,	}		return		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch accounts"})


        Color:    req.Color,        Currency: req.Currency,        Type:     req.Type,	if err != nil {	accounts, err := h.accountRepo.GetAccounts(c.Request.Context(), userID)


        Name:     req.Name,    account := &models.Account{        }        return        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})    if err := c.ShouldBindJSON(&req); err != nil {


	userID := getUserIDFromContext(c)func (h *Handler) GetAccounts(c *gin.Context) {
// Account handlers}    c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})

        }        Color           string  `json:"color" binding:"required,hexcolor"`        Currency        string  `json:"currency" binding:"required"`

        }        return
        InitialBalance  float64 `json:"initial_balance"`        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
        Type            string  `json:"type" binding:"required"`    if err := h.categoryRepo.Delete(uint(categoryID)); err != nil {    // Delete category        }        return        c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete category that is in use by transactions"})    if inUse {        }        return        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check category usage"})    if err != nil {    inUse, err := h.transactionRepo.IsCategoryInUse(uint(categoryID))

        Name            string  `json:"name" binding:"required"`



    var req struct {        userID := getUserIDFromContext(c)func (h *Handler) CreateAccount(c *gin.Context) {




// CreateAccount handles creating a new account}    c.JSON(http.StatusOK, accountsWithBalances)

        }        }            "last_updated": account.UpdatedAt,            "created_at":   account.CreatedAt,
    // Check if category is in use            "is_active":    account.IsActive,            "currency":     account.Currency,

        }        return        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})    if category.UserID != userID {


            "balance":      balance,            "color":        account.Color,            "type":         account.Type,            "name":         account.Name,
        }        return        c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})    if err != nil {    category, err := h.categoryRepo.GetByID(uint(categoryID))


            "id":           account.ID,        accountsWithBalances[i] = gin.H{

                }            balance = 0        if err != nil {        balance, err := h.transactionRepo.GetAccountBalance(account.ID)

    // Check if category exists and belongs to user        }        return        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})    if err != nil {


    for i, account := range accounts {    accountsWithBalances := make([]gin.H, len(accounts))
    // Calculate balances for each account        }        return        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accounts"})    if err != nil {    accounts, err := h.accountRepo.GetByUserID(userID)


    categoryID, err := strconv.ParseUint(c.Param("id"), 10, 64)    userID := getUserIDFromContext(c)

func (h *Handler) DeleteCategory(c *gin.Context) {// DeleteCategory handles deleting a category        userID := getUserIDFromContext(c)

}    c.JSON(http.StatusOK, category)        }// GetAccounts handles retrieving accounts
func (h *Handler) GetAccounts(c *gin.Context) {