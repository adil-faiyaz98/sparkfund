package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	// Common validation patterns
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};:'",.<>/?]{8,}$`)
	phoneRegex    = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func init() {
	validate = validator.New()

	// Register custom validation tags
	validate.RegisterValidation("email", validateEmail)
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("currency", validateCurrency)

	// Register custom type functions
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	var messages []string
	for _, err := range v {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// Validate validates a struct using tags
func Validate(v interface{}) error {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	var errors ValidationErrors

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		message := getErrorMessage(tag, field, param)
		errors = append(errors, ValidationError{
			Field:   field,
			Message: message,
		})
	}

	return errors
}

// Custom validation functions
func validateEmail(fl validator.FieldLevel) bool {
	return emailRegex.MatchString(fl.Field().String())
}

func validatePassword(fl validator.FieldLevel) bool {
	return passwordRegex.MatchString(fl.Field().String())
}

func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

func validateCurrency(fl validator.FieldLevel) bool {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "AUD", "CAD", "CHF", "CNY", "INR"}
	value := fl.Field().String()
	for _, currency := range currencies {
		if value == currency {
			return true
		}
	}
	return false
}

// getErrorMessage returns a user-friendly error message for validation errors
func getErrorMessage(tag string, field string, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, param)
	case "password":
		return fmt.Sprintf("%s must be at least 8 characters long and contain only letters, numbers, and special characters", field)
	case "phone":
		return fmt.Sprintf("%s must be a valid phone number", field)
	case "currency":
		return fmt.Sprintf("%s must be a valid currency code", field)
	default:
		return fmt.Sprintf("%s failed validation for tag %s", field, tag)
	}
}

// Example usage:
// type User struct {
//     Email    string `json:"email" validate:"required,email"`
//     Password string `json:"password" validate:"required,password"`
//     Phone    string `json:"phone" validate:"required,phone"`
// }
//
// user := User{...}
// if err := validation.Validate(user); err != nil {
//     // Handle validation errors
// }
