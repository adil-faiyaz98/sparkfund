package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/pkg/errors"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error represents an error response
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Success sends a success response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// Created sends a created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a no content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequest sends a bad request response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &Error{
			Code:    http.StatusBadRequest,
			Message: message,
		},
	})
}

// Unauthorized sends an unauthorized response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Error: &Error{
			Code:    http.StatusUnauthorized,
			Message: message,
		},
	})
}

// Forbidden sends a forbidden response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Error: &Error{
			Code:    http.StatusForbidden,
			Message: message,
		},
	})
}

// NotFound sends a not found response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error: &Error{
			Code:    http.StatusNotFound,
			Message: message,
		},
	})
}

// Conflict sends a conflict response
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success: false,
		Error: &Error{
			Code:    http.StatusConflict,
			Message: message,
		},
	})
}

// InternalServerError sends an internal server error response
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &Error{
			Code:    http.StatusInternalServerError,
			Message: message,
		},
	})
}

// ServiceUnavailable sends a service unavailable response
func ServiceUnavailable(c *gin.Context, message string) {
	c.JSON(http.StatusServiceUnavailable, Response{
		Success: false,
		Error: &Error{
			Code:    http.StatusServiceUnavailable,
			Message: message,
		},
	})
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Code, Response{
			Success: false,
			Error: &Error{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Error(),
			},
		})
		return
	}

	InternalServerError(c, "Internal server error")
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, message string, details string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: &Error{
			Code:    http.StatusBadRequest,
			Message: message,
			Details: details,
		},
	})
}

// Pagination represents pagination information
type Pagination struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, data interface{}, pagination Pagination) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: PaginatedResponse{
			Data:       data,
			Pagination: pagination,
		},
	})
}

// File sends a file response
func File(c *gin.Context, filepath string) {
	c.File(filepath)
}

// Download sends a file download response
func Download(c *gin.Context, filepath, filename string) {
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.File(filepath)
}

// Stream sends a stream response
func Stream(c *gin.Context, contentType string, reader interface{}) {
	c.Header("Content-Type", contentType)
	c.Stream(func(w io.Writer) bool {
		// TODO: Implement streaming
		return false
	})
}

// Redirect sends a redirect response
func Redirect(c *gin.Context, location string) {
	c.Redirect(http.StatusFound, location)
}

// XML sends an XML response
func XML(c *gin.Context, data interface{}) {
	c.XML(http.StatusOK, data)
}

// YAML sends a YAML response
func YAML(c *gin.Context, data interface{}) {
	c.YAML(http.StatusOK, data)
}

// Text sends a text response
func Text(c *gin.Context, format string, values ...interface{}) {
	c.String(http.StatusOK, format, values...)
}

// HTML sends an HTML response
func HTML(c *gin.Context, name string, data interface{}) {
	c.HTML(http.StatusOK, name, data)
} 