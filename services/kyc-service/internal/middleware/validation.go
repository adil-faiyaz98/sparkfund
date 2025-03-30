package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidationConfig holds validation middleware configuration
type ValidationConfig struct {
	MaxRequestSize    int64
	AllowedFileTypes  []string
	MaxFileSize       int64
	SanitizeInput     bool
	ValidateInput     bool
}

// DefaultValidationConfig returns default validation configuration
func DefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		MaxRequestSize:    10 << 20, // 10MB
		AllowedFileTypes:  []string{"application/pdf", "image/jpeg", "image/png"},
		MaxFileSize:       5 << 20,  // 5MB
		SanitizeInput:     true,
		ValidateInput:     true,
	}
}

// ValidationMiddleware provides input validation and sanitization
func ValidationMiddleware(config ValidationConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set max request size
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxRequestSize)

		// Validate content type
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "multipart/form-data") {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported media type"})
			c.Abort()
			return
		}

		// Handle file uploads
		if strings.Contains(contentType, "multipart/form-data") {
			file, err := c.FormFile("file")
			if err == nil && file != nil {
				// Validate file size
				if file.Size > config.MaxFileSize {
					c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
					c.Abort()
					return
				}

				// Validate file type
				fileType := file.Header.Get("Content-Type")
				allowed := false
				for _, t := range config.AllowedFileTypes {
					if t == fileType {
						allowed = true
						break
					}
				}
				if !allowed {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
					c.Abort()
					return
				}
			}
		}

		// Sanitize input
		if config.SanitizeInput {
			sanitizeInput(c)
		}

		// Validate input
		if config.ValidateInput {
			validateInput(c)
		}

		c.Next()
	}
}

// sanitizeInput sanitizes request input
func sanitizeInput(c *gin.Context) {
	// Sanitize query parameters
	for key, values := range c.Request.URL.Query() {
		for i, value := range values {
			values[i] = sanitizeString(value)
		}
		c.Request.URL.RawQuery = c.Request.URL.Query().Encode()
	}

	// Sanitize form values
	if err := c.Request.ParseForm(); err == nil {
		for key, values := range c.Request.Form {
			for i, value := range values {
				c.Request.Form[key][i] = sanitizeString(value)
			}
		}
	}
}

// validateInput validates request input
func validateInput(c *gin.Context) {
	validate := validator.New()

	// Validate query parameters
	for key, values := range c.Request.URL.Query() {
		for _, value := range values {
			if err := validate.Var(value, "required"); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
				c.Abort()
				return
			}
		}
	}

	// Validate form values
	if err := c.Request.ParseForm(); err == nil {
		for key, values := range c.Request.Form {
			for _, value := range values {
				if err := validate.Var(value, "required"); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
					c.Abort()
					return
				}
			}
		}
	}
}

// sanitizeString sanitizes a string input
func sanitizeString(input string) string {
	// Remove HTML tags
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")

	// Remove SQL injection attempts
	input = strings.ReplaceAll(input, "'", "''")
	input = strings.ReplaceAll(input, ";", "")
	input = strings.ReplaceAll(input, "--", "")

	// Remove XSS attempts
	input = strings.ReplaceAll(input, "javascript:", "")
	input = strings.ReplaceAll(input, "onerror=", "")
	input = strings.ReplaceAll(input, "onload=", "")

	return input
} 