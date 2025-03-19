package middleware

import (
	"net/http"
	"time"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/errors"
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// ErrorHandler middleware handles errors and returns appropriate responses
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var response ErrorResponse

			switch {
			case errors.IsValidationError(err):
				response = ErrorResponse{
					Code:    errors.ValidationError,
					Message: err.Error(),
					Status:  http.StatusBadRequest,
				}
			case errors.IsDatabaseError(err):
				response = ErrorResponse{
					Code:    errors.DatabaseError,
					Message: "Internal server error",
					Status:  http.StatusInternalServerError,
				}
				logger.Error("Database error", zap.Error(err))
			case errors.IsKafkaError(err):
				response = ErrorResponse{
					Code:    errors.KafkaError,
					Message: "Internal server error",
					Status:  http.StatusInternalServerError,
				}
				logger.Error("Kafka error", zap.Error(err))
			case errors.IsNotFoundError(err):
				response = ErrorResponse{
					Code:    errors.NotFoundError,
					Message: err.Error(),
					Status:  http.StatusNotFound,
				}
			default:
				response = ErrorResponse{
					Code:    "INTERNAL_ERROR",
					Message: "Internal server error",
					Status:  http.StatusInternalServerError,
				}
				logger.Error("Unhandled error", zap.Error(err))
			}

			c.JSON(response.Status, response)
			c.Abort()
		}
	}
}

// Recovery middleware recovers from panics
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", zap.Any("error", err))
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    "PANIC",
					Message: "Internal server error",
					Status:  http.StatusInternalServerError,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// RequestLogger middleware logs request details
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		logger.Info("Request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

// RateLimiter middleware implements a simple rate limiter
func RateLimiter(limit rate.Limit, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(limit, burst)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Code:    "RATE_LIMIT_EXCEEDED",
				Message: "Too many requests",
				Status:  http.StatusTooManyRequests,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Auth middleware validates authentication
func Auth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Next()
	}
}

// RequireRole middleware checks if the user has the required role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.Error(errors.ErrForbidden)
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			c.Error(errors.ErrForbidden)
			c.Abort()
			return
		}

		hasRole := false
		for _, r := range userRoles {
			if r == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.Error(errors.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
