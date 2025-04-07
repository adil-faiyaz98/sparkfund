package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/service"
)

// AuthController handles authentication-related endpoints
type AuthController struct {
	authService *service.AuthService
	logger      *logrus.Logger
}

// NewAuthController creates a new authentication controller
func NewAuthController(authService *service.AuthService, logger *logrus.Logger) *AuthController {
	return &AuthController{
		authService: authService,
		logger:      logger,
	}
}

// Login godoc
// @Summary Login user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: err.Error(),
		})
		return
	}

	// Get device info from request
	deviceInfo := model.DeviceInfo{
		IPAddress:    ctx.ClientIP(),
		UserAgent:    ctx.GetHeader("User-Agent"),
		DeviceType:   req.DeviceInfo.DeviceType,
		OS:           req.DeviceInfo.OS,
		Browser:      req.DeviceInfo.Browser,
		Location:     req.DeviceInfo.Location,
		CapturedTime: time.Now(),
	}

	// Login
	loginReq := model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
		MFACode:  req.MFACode,
	}

	response, err := c.authService.Login(ctx, loginReq, deviceInfo)
	if err != nil {
		c.logger.WithError(err).WithField("email", req.Email).Error("Login failed")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Login failed",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, LoginResponse{
		Token:        response.Token,
		RefreshToken: response.RefreshToken,
		ExpiresAt:    response.ExpiresAt,
		User: UserResponse{
			ID:        response.User.ID.String(),
			Email:     response.User.Email,
			FirstName: response.User.FirstName,
			LastName:  response.User.LastName,
			Role:      response.User.Role,
		},
		MFARequired: response.MFARequired,
	})
}

// RefreshToken godoc
// @Summary Refresh token
// @Description Refresh an authentication token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} LoginResponse "Token refreshed"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: err.Error(),
		})
		return
	}

	// Refresh token
	refreshReq := model.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	response, err := c.authService.RefreshToken(ctx, refreshReq)
	if err != nil {
		c.logger.WithError(err).Error("Token refresh failed")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Token refresh failed",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, LoginResponse{
		Token:        response.Token,
		RefreshToken: response.RefreshToken,
		ExpiresAt:    response.ExpiresAt,
		User: UserResponse{
			ID:        response.User.ID.String(),
			Email:     response.User.Email,
			FirstName: response.User.FirstName,
			LastName:  response.User.LastName,
			Role:      response.User.Role,
		},
		MFARequired: response.MFARequired,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Logout a user and invalidate their refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Logout request"
// @Success 200 {object} MessageResponse "Logout successful"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	var req LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: err.Error(),
		})
		return
	}

	// Logout
	err := c.authService.Logout(ctx, req.RefreshToken)
	if err != nil {
		c.logger.WithError(err).Error("Logout failed")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Logout failed",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, MessageResponse{
		Message: "Logout successful",
	})
}

// SetupMFA godoc
// @Summary Setup MFA
// @Description Setup multi-factor authentication for a user
// @Tags auth
// @Produce json
// @Success 200 {object} MFASetupResponse "MFA setup successful"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/mfa/setup [post]
func (c *AuthController) SetupMFA(ctx *gin.Context) {
	// Get user ID from JWT
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Details: err.Error(),
		})
		return
	}

	// Setup MFA
	response, err := c.authService.SetupMFA(ctx, userID)
	if err != nil {
		c.logger.WithError(err).WithField("user_id", userID).Error("MFA setup failed")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "MFA setup failed",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, MFASetupResponse{
		Secret:        response.Secret,
		QRCodeURL:     response.QRCodeURL,
		RecoveryCodes: response.RecoveryCodes,
	})
}

// VerifyMFA godoc
// @Summary Verify MFA
// @Description Verify a multi-factor authentication code
// @Tags auth
// @Accept json
// @Produce json
// @Param request body MFAVerifyRequest true "MFA verification request"
// @Success 200 {object} MessageResponse "MFA verification successful"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/mfa/verify [post]
func (c *AuthController) VerifyMFA(ctx *gin.Context) {
	var req MFAVerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
			Details: err.Error(),
		})
		return
	}

	// Get user ID from JWT
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Details: err.Error(),
		})
		return
	}

	// Verify MFA
	mfaReq := model.MFAVerifyRequest{
		Code: req.Code,
	}

	err = c.authService.VerifyMFA(ctx, userID, mfaReq)
	if err != nil {
		c.logger.WithError(err).WithField("user_id", userID).Error("MFA verification failed")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "MFA verification failed",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, MessageResponse{
		Message: "MFA verification successful",
	})
}

// DisableMFA godoc
// @Summary Disable MFA
// @Description Disable multi-factor authentication for a user
// @Tags auth
// @Produce json
// @Success 200 {object} MessageResponse "MFA disabled"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/mfa/disable [post]
func (c *AuthController) DisableMFA(ctx *gin.Context) {
	// Get user ID from JWT
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
			Details: err.Error(),
		})
		return
	}

	// Disable MFA
	err = c.authService.DisableMFA(ctx, userID)
	if err != nil {
		c.logger.WithError(err).WithField("user_id", userID).Error("MFA disabling failed")
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "MFA disabling failed",
			Details: err.Error(),
		})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, MessageResponse{
		Message: "MFA disabled",
	})
}

// RegisterRoutes registers the authentication controller routes
func (c *AuthController) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", c.Login)
		auth.POST("/refresh", c.RefreshToken)
		auth.POST("/logout", c.Logout)

		// MFA routes (require authentication)
		mfa := auth.Group("/mfa")
		mfa.Use(AuthMiddleware(c.authService))
		{
			mfa.POST("/setup", c.SetupMFA)
			mfa.POST("/verify", c.VerifyMFA)
			mfa.POST("/disable", c.DisableMFA)
		}
	}
}

// AuthMiddleware is a middleware that validates JWT tokens
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get token from Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Details: "Missing Authorization header",
			})
			return
		}

		// Extract token
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Details: err.Error(),
			})
			return
		}

		// Set claims in context
		ctx.Set("user_id", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)
		ctx.Set("mfa_passed", claims.MFAPassed)

		ctx.Next()
	}
}

// getUserIDFromContext gets the user ID from the context
func getUserIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID in context")
	}

	return userID, nil
}

// Request and response types

// LoginRequest represents a login request
type LoginRequest struct {
	Email      string            `json:"email" binding:"required,email"`
	Password   string            `json:"password" binding:"required,min=8"`
	MFACode    string            `json:"mfa_code,omitempty"`
	DeviceInfo DeviceInfoRequest `json:"device_info"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	User         UserResponse `json:"user"`
	MFARequired  bool         `json:"mfa_required,omitempty"`
}

// UserResponse represents a user in a response
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest represents a logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// MFASetupResponse represents an MFA setup response
type MFASetupResponse struct {
	Secret        string   `json:"secret"`
	QRCodeURL     string   `json:"qr_code_url"`
	RecoveryCodes []string `json:"recovery_codes"`
}

// MFAVerifyRequest represents an MFA verification request
type MFAVerifyRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}
