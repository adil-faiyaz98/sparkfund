package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sparkfund/kyc-service/internal/security"
)

// Handler struct to manage dependencies
type Handler struct {
	authConfig       security.AuthConfig
	validationConfig security.ValidationConfig
	encryptionConfig security.EncryptionConfig
	auditLogger      *security.AuditLogger
	validationRules  *security.ValidationRules
	financialRules   *security.FinancialValidationRules
}

// NewHandler creates a new handler instance with dependencies
func NewHandler(
	authConfig security.AuthConfig,
	validationConfig security.ValidationConfig,
	encryptionConfig security.EncryptionConfig,
	auditLogger *security.AuditLogger,
) *Handler {
	return &Handler{
		authConfig:       authConfig,
		validationConfig: validationConfig,
		encryptionConfig: encryptionConfig,
		auditLogger:      auditLogger,
		validationRules:  security.DefaultValidationRules(),
		financialRules:   security.DefaultFinancialRules(),
	}
}

// Handler functions
func (h *Handler) HandleLogin(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		h.auditLogger.LogSecurityEvent("login_attempt", "", map[string]interface{}{
			"error": "invalid_request_format",
			"ip":    c.ClientIP(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate input with enhanced rules
	if err := h.validationRules.ValidateUsername(loginRequest.Username); err != nil {
		h.auditLogger.LogSecurityEvent("login_attempt", "", map[string]interface{}{
			"error":    "invalid_username",
			"username": loginRequest.Username,
			"ip":       c.ClientIP(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Rate limiting check with IP-based blocking
	if !h.checkRateLimit(c.ClientIP(), "login") {
		h.auditLogger.LogSecurityEvent("rate_limit_exceeded", "", map[string]interface{}{
			"ip": c.ClientIP(),
		})
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many login attempts"})
		return
	}

	// Hash password for comparison
	hashedPassword := security.HashPassword(loginRequest.Password)

	// TODO: Verify credentials against database
	// For now, we'll just generate tokens
	userID := "user123" // This should come from the database
	roles := []string{"user"}

	// Generate tokens with enhanced security
	token, err := security.GenerateToken(userID, loginRequest.Username, roles, h.authConfig)
	if err != nil {
		h.auditLogger.LogSecurityEvent("token_generation_failed", userID, map[string]interface{}{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := security.GenerateRefreshToken(userID, h.authConfig)
	if err != nil {
		h.auditLogger.LogSecurityEvent("refresh_token_generation_failed", userID, map[string]interface{}{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Log successful authentication event with enhanced details
	h.auditLogger.LogAuthenticationEvent(userID, true, map[string]interface{}{
		"username":   loginRequest.Username,
		"ip":         c.ClientIP(),
		"user_agent": c.GetHeader("User-Agent"),
		"roles":      roles,
		"login_time": time.Now(),
		"session_id": c.GetString("session_id"),
		"device_id":  c.GetHeader("X-Device-ID"),
		"location":   c.GetHeader("X-Location"),
	})

	// Set secure cookie for refresh token with enhanced security
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(h.authConfig.RefreshTokenExpiry.Seconds()),
		"/",
		"",
		true,     // Secure
		true,     // HTTP only
		"Strict", // SameSite
	)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func HandleRefreshToken(c *gin.Context) {
	var refreshRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate refresh token
	claims := &security.Claims{}
	token, err := jwt.ParseWithClaims(refreshRequest.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(authConfig.SecretKey), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new tokens
	newToken, err := security.GenerateToken(claims.UserID, claims.Username, claims.Roles, authConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	newRefreshToken, err := security.GenerateRefreshToken(claims.UserID, authConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Log token refresh event
	auditLogger.LogSecurityEvent("token_refresh", claims.UserID, map[string]interface{}{
		"ip": c.ClientIP(),
	})

	c.JSON(http.StatusOK, gin.H{
		"token":         newToken,
		"refresh_token": newRefreshToken,
	})
}

func (h *Handler) HandleRegister(c *gin.Context) {
	var registerRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		h.auditLogger.LogSecurityEvent("registration_attempt", "", map[string]interface{}{
			"error": "invalid_request_format",
			"ip":    c.ClientIP(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate input with enhanced rules
	if err := h.validationRules.ValidateUsername(registerRequest.Username); err != nil {
		h.auditLogger.LogSecurityEvent("registration_attempt", "", map[string]interface{}{
			"error":    "invalid_username",
			"username": registerRequest.Username,
			"ip":       c.ClientIP(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validationRules.ValidatePassword(registerRequest.Password); err != nil {
		h.auditLogger.LogSecurityEvent("registration_attempt", "", map[string]interface{}{
			"error":    "invalid_password",
			"username": registerRequest.Username,
			"ip":       c.ClientIP(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validationRules.ValidateEmail(registerRequest.Email); err != nil {
		h.auditLogger.LogSecurityEvent("registration_attempt", "", map[string]interface{}{
			"error":    "invalid_email",
			"username": registerRequest.Username,
			"ip":       c.ClientIP(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword := security.HashPassword(registerRequest.Password)

	// TODO: Store user in database
	// For now, we'll just return success
	userID := "user123" // This should be generated

	// Log registration event with enhanced details
	h.auditLogger.LogSecurityEvent("user_registration", userID, map[string]interface{}{
		"username":      registerRequest.Username,
		"email":         security.MaskSensitiveData(registerRequest.Email),
		"ip":            c.ClientIP(),
		"user_agent":    c.GetHeader("User-Agent"),
		"register_time": time.Now(),
	})

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *Handler) HandleCreateKYC(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		h.auditLogger.LogSecurityEvent("kyc_creation_attempt", "", map[string]interface{}{
			"error": "user_not_authenticated",
			"ip":    c.ClientIP(),
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Validate file upload
	file, err := c.FormFile("document")
	if err != nil {
		h.auditLogger.LogSecurityEvent("kyc_creation_attempt", userID.(string), map[string]interface{}{
			"error": "no_file_uploaded",
			"ip":    c.ClientIP(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Get document type from form
	docType := c.PostForm("document_type")
	if docType == "" {
		h.auditLogger.LogSecurityEvent("kyc_creation_attempt", userID.(string), map[string]interface{}{
			"error":     "missing_document_type",
			"file_name": file.Filename,
			"ip":        c.ClientIP(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Document type is required"})
		return
	}

	// Validate document according to financial rules
	if err := h.financialRules.ValidateDocument(docType, nil, time.Now()); err != nil {
		h.auditLogger.LogSecurityEvent("kyc_creation_attempt", userID.(string), map[string]interface{}{
			"error":     "invalid_document",
			"file_name": file.Filename,
			"doc_type":  docType,
			"ip":        c.ClientIP(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Read file content
	src, err := file.Open()
	if err != nil {
		h.auditLogger.LogSecurityEvent("kyc_creation_attempt", userID.(string), map[string]interface{}{
			"error":     "file_read_error",
			"file_name": file.Filename,
			"ip":        c.ClientIP(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer src.Close()

	// Read file content into buffer
	buf := make([]byte, file.Size)
	if _, err := src.Read(buf); err != nil {
		h.auditLogger.LogSecurityEvent("kyc_creation_attempt", userID.(string), map[string]interface{}{
			"error":     "file_read_error",
			"file_name": file.Filename,
			"ip":        c.ClientIP(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Hash file for integrity check
	fileHash := security.HashFile(buf)

	// Encrypt file content
	encryptedData, err := h.encryptSensitiveData(buf)
	if err != nil {
		h.auditLogger.LogSecurityEvent("kyc_creation_attempt", userID.(string), map[string]interface{}{
			"error":     "encryption_error",
			"file_name": file.Filename,
			"ip":        c.ClientIP(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt file"})
		return
	}

	// TODO: Store encrypted data and metadata in database

	// Log data access event with enhanced details
	h.auditLogger.LogDataAccessEvent(userID.(string), "kyc_document", "create", map[string]interface{}{
		"file_name":  file.Filename,
		"file_size":  file.Size,
		"file_hash":  fileHash,
		"doc_type":   docType,
		"ip":         c.ClientIP(),
		"user_agent": c.GetHeader("User-Agent"),
		"timestamp":  time.Now(),
		"session_id": c.GetString("session_id"),
		"device_id":  c.GetHeader("X-Device-ID"),
		"location":   c.GetHeader("X-Location"),
		"risk_level": h.calculateRiskLevel(docType),
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "KYC document uploaded successfully",
		"hash":    fileHash,
	})
}

func HandleGetKYC(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	requestedUserID := c.Param("userID")

	// Check if user has permission to access this KYC
	if userID.(string) != requestedUserID {
		roles, _ := c.Get("roles")
		hasAdminRole := false
		for _, role := range roles.([]string) {
			if role == "admin" {
				hasAdminRole = true
				break
			}
		}
		if !hasAdminRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
	}

	// TODO: Retrieve KYC data from database
	// For now, return mock data
	kycData := map[string]interface{}{
		"user_id": requestedUserID,
		"status":  "pending",
	}

	// Log data access event
	auditLogger.LogDataAccessEvent(userID.(string), "kyc_data", "read", map[string]interface{}{
		"requested_user_id": requestedUserID,
	})

	c.JSON(http.StatusOK, kycData)
}

func HandleUpdateKYC(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var updateRequest struct {
		UserID      string                 `json:"user_id" binding:"required"`
		Data        map[string]interface{} `json:"data" binding:"required"`
		DocumentIDs []string               `json:"document_ids"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate input
	if err := security.ValidateString(updateRequest.UserID, validationConfig); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Check permissions
	if userID.(string) != updateRequest.UserID {
		roles, _ := c.Get("roles")
		hasAdminRole := false
		for _, role := range roles.([]string) {
			if role == "admin" {
				hasAdminRole = true
				break
			}
		}
		if !hasAdminRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
	}

	// TODO: Update KYC data in database

	// Log data access event
	auditLogger.LogDataAccessEvent(userID.(string), "kyc_data", "update", map[string]interface{}{
		"updated_user_id": updateRequest.UserID,
		"document_ids":    updateRequest.DocumentIDs,
	})

	c.JSON(http.StatusOK, gin.H{"message": "KYC data updated successfully"})
}

func HandleUpdateKYCStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has admin role
	roles, _ := c.Get("roles")
	hasAdminRole := false
	for _, role := range roles.([]string) {
		if role == "admin" {
			hasAdminRole = true
			break
		}
	}
	if !hasAdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var statusRequest struct {
		Status string `json:"status" binding:"required,oneof=pending approved rejected"`
	}

	if err := c.ShouldBindJSON(&statusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	requestedUserID := c.Param("userID")

	// TODO: Update KYC status in database

	// Log configuration change event
	auditLogger.LogConfigurationChangeEvent(
		userID.(string),
		"kyc_status",
		"pending",
		statusRequest.Status,
		map[string]interface{}{
			"user_id": requestedUserID,
		},
	)

	c.JSON(http.StatusOK, gin.H{"message": "KYC status updated successfully"})
}

func HandleListKYC(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has admin role
	roles, _ := c.Get("roles")
	hasAdminRole := false
	for _, role := range roles.([]string) {
		if role == "admin" {
			hasAdminRole = true
			break
		}
	}
	if !hasAdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	// TODO: Retrieve KYC list from database
	// For now, return mock data
	kycList := []map[string]interface{}{
		{
			"user_id": "user123",
			"status":  "pending",
		},
	}

	// Log data access event
	auditLogger.LogDataAccessEvent(userID.(string), "kyc_list", "read", map[string]interface{}{
		"list_size": len(kycList),
	})

	c.JSON(http.StatusOK, kycList)
}

func HandleDeleteKYC(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has admin role
	roles, _ := c.Get("roles")
	hasAdminRole := false
	for _, role := range roles.([]string) {
		if role == "admin" {
			hasAdminRole = true
			break
		}
	}
	if !hasAdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	requestedUserID := c.Param("userID")

	// TODO: Delete KYC data from database

	// Log data access event
	auditLogger.LogDataAccessEvent(userID.(string), "kyc_data", "delete", map[string]interface{}{
		"deleted_user_id": requestedUserID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "KYC data deleted successfully"})
}

func HandleGetAuditLogs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has admin role
	roles, _ := c.Get("roles")
	hasAdminRole := false
	for _, role := range roles.([]string) {
		if role == "admin" {
			hasAdminRole = true
			break
		}
	}
	if !hasAdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	// TODO: Implement audit log retrieval
	// For now, return mock data
	logs := []map[string]interface{}{
		{
			"timestamp": time.Now(),
			"event":     "user_login",
			"user_id":   "user123",
		},
	}

	// Log data access event
	auditLogger.LogDataAccessEvent(userID.(string), "audit_logs", "read", map[string]interface{}{
		"log_count": len(logs),
	})

	c.JSON(http.StatusOK, logs)
}

func HandleGetMetrics(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has admin role
	roles, _ := c.Get("roles")
	hasAdminRole := false
	for _, role := range roles.([]string) {
		if role == "admin" {
			hasAdminRole = true
			break
		}
	}
	if !hasAdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	// TODO: Implement metrics retrieval
	// For now, return mock data
	metrics := map[string]interface{}{
		"total_users": 100,
		"active_kyc":  50,
	}

	// Log data access event
	auditLogger.LogDataAccessEvent(userID.(string), "metrics", "read", nil)

	c.JSON(http.StatusOK, metrics)
}

func HandleUpdateConfig(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has admin role
	roles, _ := c.Get("roles")
	hasAdminRole := false
	for _, role := range roles.([]string) {
		if role == "admin" {
			hasAdminRole = true
			break
		}
	}
	if !hasAdminRole {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		return
	}

	var configRequest struct {
		Key   string      `json:"key" binding:"required"`
		Value interface{} `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&configRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// TODO: Update configuration in database

	// Log configuration change event
	auditLogger.LogConfigurationChangeEvent(
		userID.(string),
		configRequest.Key,
		nil, // Old value
		configRequest.Value,
		nil,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
}

func HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func HandleMetrics(c *gin.Context) {
	// TODO: Implement Prometheus metrics endpoint
	c.JSON(http.StatusOK, gin.H{"message": "Metrics endpoint"})
}

// Helper functions
func (h *Handler) checkRateLimit(ip, action string) bool {
	// TODO: Implement rate limiting using Redis
	return true
}

func (h *Handler) validateUserPermissions(c *gin.Context, requiredRoles []string) bool {
	userID, exists := c.Get("user_id")
	if !exists {
		return false
	}

	roles, exists := c.Get("roles")
	if !exists {
		return false
	}

	userRoles := roles.([]string)
	for _, requiredRole := range requiredRoles {
		for _, userRole := range userRoles {
			if userRole == requiredRole {
				return true
			}
		}
	}

	h.auditLogger.LogAuthorizationEvent(userID.(string), "permission_denied", map[string]interface{}{
		"required_roles": requiredRoles,
		"user_roles":     userRoles,
		"ip":             c.ClientIP(),
		"session_id":     c.GetString("session_id"),
		"device_id":      c.GetHeader("X-Device-ID"),
		"location":       c.GetHeader("X-Location"),
	})

	return false
}

func (h *Handler) validateAndSanitizeInput(input interface{}) error {
	// TODO: Implement input validation and sanitization
	return nil
}

func (h *Handler) encryptSensitiveData(data []byte) ([]byte, error) {
	return security.EncryptAES(data, h.encryptionConfig)
}

func (h *Handler) decryptSensitiveData(data []byte) ([]byte, error) {
	return security.DecryptAES(data, h.encryptionConfig)
}

func (h *Handler) calculateRiskLevel(docType string) string {
	// TODO: Implement risk level calculation based on document type and other factors
	return "medium"
}
