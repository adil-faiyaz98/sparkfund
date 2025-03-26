package validation

import (
	"fmt"
	"net/http"
	"regexp"

	"sparkfund/email-service/internal/errors"
	"sparkfund/email-service/internal/models"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// validateEmail validates an email address format
func validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.NewError(
			http.StatusBadRequest,
			fmt.Sprintf("invalid email address: %s", email),
		)
	}
	return nil
}

// ValidateEmailRequest validates a SendEmailRequest
func ValidateEmailRequest(req *models.SendEmailRequest, fromAddress string) error {
	if err := validateEmail(fromAddress); err != nil {
		return err
	}

	for _, to := range req.To {
		if err := validateEmail(string(to)); err != nil {
			return err
		}
	}

	for _, cc := range req.Cc {
		if err := validateEmail(string(cc)); err != nil {
			return err
		}
	}

	for _, bcc := range req.Bcc {
		if err := validateEmail(string(bcc)); err != nil {
			return err
		}
	}

	return models.ValidateStruct(req)
}

// ValidateTemplate validates an email template
func ValidateTemplate(template *models.Template) error {
	return models.ValidateStruct(template)
}

// ValidateEmailLog validates email log data
func ValidateEmailLog(log *models.EmailLog) error {
	return models.ValidateStruct(log)
}