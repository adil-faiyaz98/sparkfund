package validator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sparkfund/pkg/errors"
)

var (
	validate *validator.Validate

	// Regular expressions
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	uuidRegex     = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	dateRegex     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	amountRegex   = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	percentRegex  = regexp.MustCompile(`^(?:100|[1-9]?\d)(?:\.\d{1,2})?$`)
	urlRegex      = regexp.MustCompile(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([/\w \.-]*)*\/?$`)
	passwordRegex = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`)
)

// Init initializes the validator
func Init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using tags
func ValidateStruct(s interface{}) error {
	if validate == nil {
		Init()
	}
	if err := validate.Struct(s); err != nil {
		return errors.ErrBadRequest(err)
	}
	return nil
}

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.ErrBadRequest(fmt.Errorf("invalid email format"))
	}
	return nil
}

// ValidatePhone validates a phone number
func ValidatePhone(phone string) error {
	if !phoneRegex.MatchString(phone) {
		return errors.ErrBadRequest(fmt.Errorf("invalid phone number format"))
	}
	return nil
}

// ValidatePassword validates a password
func ValidatePassword(password string) error {
	if !passwordRegex.MatchString(password) {
		return errors.ErrBadRequest(fmt.Errorf("password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number and one special character"))
	}
	return nil
}

// ValidateUUID validates a UUID
func ValidateUUID(uuid string) error {
	if !uuidRegex.MatchString(uuid) {
		return errors.ErrBadRequest(fmt.Errorf("invalid UUID format"))
	}
	return nil
}

// ValidateDate validates a date string
func ValidateDate(date string) error {
	if !dateRegex.MatchString(date) {
		return errors.ErrBadRequest(fmt.Errorf("invalid date format, expected YYYY-MM-DD"))
	}
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return errors.ErrBadRequest(fmt.Errorf("invalid date"))
	}
	return nil
}

// ValidateAmount validates a monetary amount
func ValidateAmount(amount string) error {
	if !amountRegex.MatchString(amount) {
		return errors.ErrBadRequest(fmt.Errorf("invalid amount format"))
	}
	return nil
}

// ValidatePercentage validates a percentage value
func ValidatePercentage(percent string) error {
	if !percentRegex.MatchString(percent) {
		return errors.ErrBadRequest(fmt.Errorf("invalid percentage format"))
	}
	return nil
}

// ValidateURL validates a URL
func ValidateURL(url string) error {
	if !urlRegex.MatchString(url) {
		return errors.ErrBadRequest(fmt.Errorf("invalid URL format"))
	}
	return nil
}

// ValidateRequired validates that a field is not empty
func ValidateRequired(value interface{}) error {
	if value == nil {
		return errors.ErrBadRequest(fmt.Errorf("field is required"))
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			return errors.ErrBadRequest(fmt.Errorf("field is required"))
		}
	case []string:
		if len(v) == 0 {
			return errors.ErrBadRequest(fmt.Errorf("field is required"))
		}
	case []interface{}:
		if len(v) == 0 {
			return errors.ErrBadRequest(fmt.Errorf("field is required"))
		}
	}
	return nil
}

// ValidateMinLength validates minimum length of a string
func ValidateMinLength(value string, min int) error {
	if len(value) < min {
		return errors.ErrBadRequest(fmt.Errorf("field must be at least %d characters long", min))
	}
	return nil
}

// ValidateMaxLength validates maximum length of a string
func ValidateMaxLength(value string, max int) error {
	if len(value) > max {
		return errors.ErrBadRequest(fmt.Errorf("field must not exceed %d characters", max))
	}
	return nil
}

// ValidateMin validates minimum value of a number
func ValidateMin(value float64, min float64) error {
	if value < min {
		return errors.ErrBadRequest(fmt.Errorf("value must be at least %f", min))
	}
	return nil
}

// ValidateMax validates maximum value of a number
func ValidateMax(value float64, max float64) error {
	if value > max {
		return errors.ErrBadRequest(fmt.Errorf("value must not exceed %f", max))
	}
	return nil
}

// ValidateRange validates that a number is within a range
func ValidateRange(value float64, min, max float64) error {
	if value < min || value > max {
		return errors.ErrBadRequest(fmt.Errorf("value must be between %f and %f", min, max))
	}
	return nil
}

// ValidateIn validates that a value is in a list of allowed values
func ValidateIn(value string, allowed []string) error {
	for _, v := range allowed {
		if v == value {
			return nil
		}
	}
	return errors.ErrBadRequest(fmt.Errorf("value must be one of: %v", allowed))
}

// ValidateNotIn validates that a value is not in a list of forbidden values
func ValidateNotIn(value string, forbidden []string) error {
	for _, v := range forbidden {
		if v == value {
			return errors.ErrBadRequest(fmt.Errorf("value must not be one of: %v", forbidden))
		}
	}
	return nil
}

// ValidateUnique validates that all values in a slice are unique
func ValidateUnique(values []string) error {
	seen := make(map[string]bool)
	for _, v := range values {
		if seen[v] {
			return errors.ErrBadRequest(fmt.Errorf("duplicate value: %s", v))
		}
		seen[v] = true
	}
	return nil
}

// ValidateRegex validates that a string matches a regular expression
func ValidateRegex(value, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return errors.ErrInternalServer(err)
	}
	if !re.MatchString(value) {
		return errors.ErrBadRequest(fmt.Errorf("value does not match pattern"))
	}
	return nil
} 