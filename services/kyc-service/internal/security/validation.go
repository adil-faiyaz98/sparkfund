package security

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

// ValidationConfig holds validation configuration
type ValidationConfig struct {
	MaxFileSize    int64
	AllowedTypes   []string
	MaxStringLength int
	MinStringLength int
}

// DefaultValidationConfig returns default validation configuration
func DefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		MaxFileSize:     10 * 1024 * 1024, // 10MB
		AllowedTypes:    []string{"image/jpeg", "image/png", "application/pdf"},
		MaxStringLength: 1000,
		MinStringLength: 1,
	}
}

// ValidateFile middleware validates file uploads
func ValidateFile(config ValidationConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			c.Abort()
			return
		}

		// Check file size
		if file.Size > config.MaxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
			c.Abort()
			return
		}

		// Check file type
		contentType := file.Header.Get("Content-Type")
		allowed := false
		for _, t := range config.AllowedTypes {
			if contentType == t {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
			c.Abort()
			return
		}

		// Check file extension
		ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, ".")+1:])
		if !isAllowedExtension(ext) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateString validates string input
func ValidateString(str string, config ValidationConfig) error {
	if len(str) > config.MaxStringLength {
		return fmt.Errorf("string too long")
	}
	if len(str) < config.MinStringLength {
		return fmt.Errorf("string too short")
	}
	return nil
}

// SanitizeHTML sanitizes HTML input
func SanitizeHTML(html string) string {
	p := bluemonday.UGCPolicy()
	return p.Sanitize(html)
}

// ValidateEmail validates email address
func ValidateEmail(email string) error {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// Check for uppercase letters
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	// Check for lowercase letters
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	// Check for numbers
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}

	// Check for special characters
	if !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidateBase64 validates base64 encoded data
func ValidateBase64(data string) error {
	_, err := base64.StdEncoding.DecodeString(data)
	return err
}

// ValidateURL validates URL format
func ValidateURL(url string) error {
	pattern := `^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`
	matched, err := regexp.MatchString(pattern, url)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("invalid URL format")
	}
	return nil
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) error {
	pattern := `^\+?[1-9]\d{1,14}$`
	matched, err := regexp.MatchString(pattern, phone)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("invalid phone number format")
	}
	return nil
}

// ValidateDate validates date format
func ValidateDate(date string) error {
	pattern := `^\d{4}-\d{2}-\d{2}$`
	matched, err := regexp.MatchString(pattern, date)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("invalid date format")
	}
	return nil
}

// isAllowedExtension checks if the file extension is allowed
func isAllowedExtension(ext string) bool {
	allowed := []string{"jpg", "jpeg", "png", "pdf"}
	for _, e := range allowed {
		if ext == e {
			return true
		}
	}
	return false
} 