package validation

import (
	"fmt"
	"regexp"
	"time"

	"github.com/sparkfund/kyc-service/internal/model"
)

var (
	postalCodeRegex     = regexp.MustCompile(`^[0-9]{5,10}$`)
	documentNumberRegex = regexp.MustCompile(`^[A-Z0-9]{5,20}$`)
)

type KYCValidator struct{}

func NewKYCValidator() *KYCValidator {
	return &KYCValidator{}
}

func (v *KYCValidator) ValidateKYCRequest(req *model.KYCRequest) error {
	if err := v.validateName(req.FirstName, "first name"); err != nil {
		return err
	}
	if err := v.validateName(req.LastName, "last name"); err != nil {
		return err
	}
	if err := v.validateDateOfBirth(req.DateOfBirth); err != nil {
		return err
	}
	if err := v.validateAddress(req.Address); err != nil {
		return err
	}
	if err := v.validateCity(req.City); err != nil {
		return err
	}
	if err := v.validateCountry(req.Country); err != nil {
		return err
	}
	if err := v.validatePostalCode(req.PostalCode); err != nil {
		return err
	}
	if err := v.validateDocumentType(req.DocumentType); err != nil {
		return err
	}
	if err := v.validateDocumentNumber(req.DocumentNumber); err != nil {
		return err
	}
	if err := v.validateDocumentImages(req.DocumentFront, req.DocumentBack, req.SelfieImage); err != nil {
		return err
	}

	return nil
}

func (v *KYCValidator) validateName(name, field string) error {
	if name == "" {
		return fmt.Errorf("%s is required", field)
	}
	if len(name) < 2 || len(name) > 50 {
		return fmt.Errorf("%s must be between 2 and 50 characters", field)
	}
	if !regexp.MustCompile(`^[a-zA-Z\s-]+$`).MatchString(name) {
		return fmt.Errorf("%s can only contain letters, spaces, and hyphens", field)
	}
	return nil
}

func (v *KYCValidator) validateDateOfBirth(dob string) error {
	if dob == "" {
		return fmt.Errorf("date of birth is required")
	}

	date, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return fmt.Errorf("invalid date of birth format. Use YYYY-MM-DD")
	}

	// Check if person is at least 18 years old
	age := time.Since(date).Hours() / 24 / 365
	if age < 18 {
		return fmt.Errorf("person must be at least 18 years old")
	}

	// Check if date is not in the future
	if date.After(time.Now()) {
		return fmt.Errorf("date of birth cannot be in the future")
	}

	return nil
}

func (v *KYCValidator) validateAddress(address string) error {
	if address == "" {
		return fmt.Errorf("address is required")
	}
	if len(address) < 5 || len(address) > 200 {
		return fmt.Errorf("address must be between 5 and 200 characters")
	}
	return nil
}

func (v *KYCValidator) validateCity(city string) error {
	if city == "" {
		return fmt.Errorf("city is required")
	}
	if len(city) < 2 || len(city) > 100 {
		return fmt.Errorf("city must be between 2 and 100 characters")
	}
	return nil
}

func (v *KYCValidator) validateCountry(country string) error {
	if country == "" {
		return fmt.Errorf("country is required")
	}
	if len(country) != 2 {
		return fmt.Errorf("country must be a 2-letter ISO code")
	}
	return nil
}

func (v *KYCValidator) validatePostalCode(postalCode string) error {
	if postalCode == "" {
		return fmt.Errorf("postal code is required")
	}
	if !postalCodeRegex.MatchString(postalCode) {
		return fmt.Errorf("invalid postal code format")
	}
	return nil
}

func (v *KYCValidator) validateDocumentType(docType string) error {
	if docType == "" {
		return fmt.Errorf("document type is required")
	}
	validTypes := map[string]bool{
		"passport":         true,
		"national_id":      true,
		"drivers_license":  true,
		"residence_permit": true,
	}
	if !validTypes[docType] {
		return fmt.Errorf("invalid document type")
	}
	return nil
}

func (v *KYCValidator) validateDocumentNumber(docNumber string) error {
	if docNumber == "" {
		return fmt.Errorf("document number is required")
	}
	if !documentNumberRegex.MatchString(docNumber) {
		return fmt.Errorf("invalid document number format")
	}
	return nil
}

func (v *KYCValidator) validateDocumentImages(front, back, selfie string) error {
	if front == "" {
		return fmt.Errorf("document front image is required")
	}
	if back == "" {
		return fmt.Errorf("document back image is required")
	}
	if selfie == "" {
		return fmt.Errorf("selfie image is required")
	}
	return nil
}
