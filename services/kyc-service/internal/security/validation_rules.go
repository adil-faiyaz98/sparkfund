package security

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationRules defines the validation rules for different types of data
type ValidationRules struct {
	Username struct {
		MinLength    int
		MaxLength    int
		AllowedChars string
		Pattern      *regexp.Regexp
	}
	Password struct {
		MinLength      int
		MaxLength      int
		RequireUpper   bool
		RequireLower   bool
		RequireNumber  bool
		RequireSpecial bool
		Pattern        *regexp.Regexp
	}
	Email struct {
		Pattern *regexp.Regexp
	}
	Document struct {
		MaxSize      int64
		AllowedTypes []string
		MaxFiles     int
	}
	KYCData struct {
		RequiredFields []string
		FieldPatterns  map[string]*regexp.Regexp
		FieldLengths   map[string]struct {
			Min int
			Max int
		}
	}
}

// DefaultValidationRules returns the default validation rules
func DefaultValidationRules() *ValidationRules {
	return &ValidationRules{
		Username: struct {
			MinLength    int
			MaxLength    int
			AllowedChars string
			Pattern      *regexp.Regexp
		}{
			MinLength:    3,
			MaxLength:    32,
			AllowedChars: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-",
			Pattern:      regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`),
		},
		Password: struct {
			MinLength      int
			MaxLength      int
			RequireUpper   bool
			RequireLower   bool
			RequireNumber  bool
			RequireSpecial bool
			Pattern        *regexp.Regexp
		}{
			MinLength:      8,
			MaxLength:      64,
			RequireUpper:   true,
			RequireLower:   true,
			RequireNumber:  true,
			RequireSpecial: true,
			Pattern:        regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,64}$`),
		},
		Email: struct {
			Pattern *regexp.Regexp
		}{
			Pattern: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		},
		Document: struct {
			MaxSize      int64
			AllowedTypes []string
			MaxFiles     int
		}{
			MaxSize:      10 * 1024 * 1024, // 10MB
			AllowedTypes: []string{".pdf", ".jpg", ".jpeg", ".png"},
			MaxFiles:     5,
		},
		KYCData: struct {
			RequiredFields []string
			FieldPatterns  map[string]*regexp.Regexp
			FieldLengths   map[string]struct {
				Min int
				Max int
			}
		}{
			RequiredFields: []string{"full_name", "date_of_birth", "address", "document_number"},
			FieldPatterns: map[string]*regexp.Regexp{
				"full_name":       regexp.MustCompile(`^[a-zA-Z\s]{2,100}$`),
				"document_number": regexp.MustCompile(`^[A-Z0-9]{5,20}$`),
				"phone":           regexp.MustCompile(`^\+?[1-9]\d{1,14}$`),
			},
			FieldLengths: map[string]struct {
				Min int
				Max int
			}{
				"address": {10, 200},
				"city":    {2, 50},
				"country": {2, 50},
			},
		},
	}
}

// ValidateUsername validates a username according to the rules
func (r *ValidationRules) ValidateUsername(username string) error {
	if len(username) < r.Username.MinLength {
		return fmt.Errorf("username must be at least %d characters long", r.Username.MinLength)
	}
	if len(username) > r.Username.MaxLength {
		return fmt.Errorf("username must not exceed %d characters", r.Username.MaxLength)
	}
	if !r.Username.Pattern.MatchString(username) {
		return fmt.Errorf("username contains invalid characters")
	}
	return nil
}

// ValidatePassword validates a password according to the rules
func (r *ValidationRules) ValidatePassword(password string) error {
	if len(password) < r.Password.MinLength {
		return fmt.Errorf("password must be at least %d characters long", r.Password.MinLength)
	}
	if len(password) > r.Password.MaxLength {
		return fmt.Errorf("password must not exceed %d characters", r.Password.MaxLength)
	}
	if r.Password.RequireUpper && !strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if r.Password.RequireLower && !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz") {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if r.Password.RequireNumber && !strings.ContainsAny(password, "0123456789") {
		return fmt.Errorf("password must contain at least one number")
	}
	if r.Password.RequireSpecial && !strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?") {
		return fmt.Errorf("password must contain at least one special character")
	}
	return nil
}

// ValidateEmail validates an email address according to the rules
func (r *ValidationRules) ValidateEmail(email string) error {
	if !r.Email.Pattern.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidateDocument validates a document according to the rules
func (r *ValidationRules) ValidateDocument(filename string, size int64) error {
	if size > r.Document.MaxSize {
		return fmt.Errorf("document size exceeds maximum allowed size of %d bytes", r.Document.MaxSize)
	}

	ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	allowed := false
	for _, t := range r.Document.AllowedTypes {
		if ext == t {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("document type not allowed. Allowed types: %v", r.Document.AllowedTypes)
	}

	return nil
}

// ValidateKYCData validates KYC data according to the rules
func (r *ValidationRules) ValidateKYCData(data map[string]interface{}) error {
	// Check required fields
	for _, field := range r.KYCData.RequiredFields {
		if _, exists := data[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate field patterns
	for field, value := range data {
		if pattern, exists := r.KYCData.FieldPatterns[field]; exists {
			strValue, ok := value.(string)
			if !ok {
				return fmt.Errorf("field %s must be a string", field)
			}
			if !pattern.MatchString(strValue) {
				return fmt.Errorf("invalid format for field: %s", field)
			}
		}

		// Validate field lengths
		if lengths, exists := r.KYCData.FieldLengths[field]; exists {
			strValue, ok := value.(string)
			if !ok {
				return fmt.Errorf("field %s must be a string", field)
			}
			if len(strValue) < lengths.Min {
				return fmt.Errorf("field %s must be at least %d characters long", field, lengths.Min)
			}
			if len(strValue) > lengths.Max {
				return fmt.Errorf("field %s must not exceed %d characters", field, lengths.Max)
			}
		}
	}

	return nil
}
