package security

import (
	"fmt"
	"regexp"
	"time"
)

// FinancialValidationRules defines validation rules specific to financial regulations
type FinancialValidationRules struct {
	Customer struct {
		MinAge           int
		RequiredFields   []string
		FieldPatterns    map[string]*regexp.Regexp
		RestrictedFields []string
		Sanitization     map[string]func(string) string
	}
	Document struct {
		RequiredTypes  []string
		MaxAge         time.Duration
		MinQuality     float64
		MaxFileSize    int64
		AllowedFormats []string
		ContentChecks  []func([]byte) error
	}
	Transaction struct {
		MaxAmount      float64
		MinAmount      float64
		RequiredFields []string
		RiskLevels     map[string]struct {
			Threshold float64
			Checks    []func(map[string]interface{}) error
		}
	}
	Compliance struct {
		WatchlistChecks       []func(string) error
		SanctionChecks        []func(string) error
		PEPChecks             []func(string) error
		RiskScoring           map[string]float64
		RequiredVerifications []string
	}
}

// DefaultFinancialRules returns default financial validation rules
func DefaultFinancialRules() *FinancialValidationRules {
	return &FinancialValidationRules{
		Customer: struct {
			MinAge           int
			RequiredFields   []string
			FieldPatterns    map[string]*regexp.Regexp
			RestrictedFields []string
			Sanitization     map[string]func(string) string
		}{
			MinAge: 18,
			RequiredFields: []string{
				"full_name",
				"date_of_birth",
				"nationality",
				"residential_address",
				"identification_number",
				"tax_id",
				"occupation",
				"source_of_funds",
				"purpose_of_account",
			},
			FieldPatterns: map[string]*regexp.Regexp{
				"full_name":             regexp.MustCompile(`^[a-zA-Z\s]{2,100}$`),
				"identification_number": regexp.MustCompile(`^[A-Z0-9]{5,20}$`),
				"tax_id":                regexp.MustCompile(`^[A-Z0-9]{10,20}$`),
				"phone":                 regexp.MustCompile(`^\+?[1-9]\d{1,14}$`),
				"email":                 regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
			},
			RestrictedFields: []string{
				"password",
				"credit_card",
				"bank_account",
				"ssn",
			},
			Sanitization: map[string]func(string) string{
				"full_name": func(s string) string {
					return SanitizeHTML(s)
				},
				"address": func(s string) string {
					return SanitizeHTML(s)
				},
			},
		},
		Document: struct {
			RequiredTypes  []string
			MaxAge         time.Duration
			MinQuality     float64
			MaxFileSize    int64
			AllowedFormats []string
			ContentChecks  []func([]byte) error
		}{
			RequiredTypes: []string{
				"identity_document",
				"proof_of_address",
				"tax_document",
			},
			MaxAge:         90 * 24 * time.Hour, // 90 days
			MinQuality:     0.8,
			MaxFileSize:    10 * 1024 * 1024, // 10MB
			AllowedFormats: []string{".pdf", ".jpg", ".jpeg", ".png"},
			ContentChecks: []func([]byte) error{
				func(data []byte) error {
					// TODO: Implement document quality check
					return nil
				},
				func(data []byte) error {
					// TODO: Implement OCR check
					return nil
				},
				func(data []byte) error {
					// TODO: Implement tampering detection
					return nil
				},
			},
		},
		Transaction: struct {
			MaxAmount      float64
			MinAmount      float64
			RequiredFields []string
			RiskLevels     map[string]struct {
				Threshold float64
				Checks    []func(map[string]interface{}) error
			}
		}{
			MaxAmount: 1000000.0,
			MinAmount: 0.01,
			RequiredFields: []string{
				"amount",
				"currency",
				"purpose",
				"source",
				"destination",
				"timestamp",
			},
			RiskLevels: map[string]struct {
				Threshold float64
				Checks    []func(map[string]interface{}) error
			}{
				"low": {
					Threshold: 10000.0,
					Checks: []func(map[string]interface{}) error{
						func(data map[string]interface{}) error {
							// TODO: Implement basic AML check
							return nil
						},
					},
				},
				"medium": {
					Threshold: 50000.0,
					Checks: []func(map[string]interface{}) error{
						func(data map[string]interface{}) error {
							// TODO: Implement enhanced AML check
							return nil
						},
						func(data map[string]interface{}) error {
							// TODO: Implement PEP check
							return nil
						},
					},
				},
				"high": {
					Threshold: 100000.0,
					Checks: []func(map[string]interface{}) error{
						func(data map[string]interface{}) error {
							// TODO: Implement comprehensive AML check
							return nil
						},
						func(data map[string]interface{}) error {
							// TODO: Implement sanctions check
							return nil
						},
						func(data map[string]interface{}) error {
							// TODO: Implement enhanced due diligence
							return nil
						},
					},
				},
			},
		},
		Compliance: struct {
			WatchlistChecks       []func(string) error
			SanctionChecks        []func(string) error
			PEPChecks             []func(string) error
			RiskScoring           map[string]float64
			RequiredVerifications []string
		}{
			WatchlistChecks: []func(string) error{
				func(name string) error {
					// TODO: Implement watchlist check
					return nil
				},
			},
			SanctionChecks: []func(string) error{
				func(name string) error {
					// TODO: Implement sanctions check
					return nil
				},
			},
			PEPChecks: []func(string) error{
				func(name string) error {
					// TODO: Implement PEP check
					return nil
				},
			},
			RiskScoring: map[string]float64{
				"country_risk":     0.3,
				"customer_risk":    0.3,
				"product_risk":     0.2,
				"transaction_risk": 0.2,
			},
			RequiredVerifications: []string{
				"identity",
				"address",
				"tax_status",
				"source_of_funds",
				"purpose_of_account",
			},
		},
	}
}

// ValidateCustomerData validates customer data according to financial regulations
func (r *FinancialValidationRules) ValidateCustomerData(data map[string]interface{}) error {
	// Check required fields
	for _, field := range r.Customer.RequiredFields {
		if _, exists := data[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate field patterns
	for field, value := range data {
		if pattern, exists := r.Customer.FieldPatterns[field]; exists {
			strValue, ok := value.(string)
			if !ok {
				return fmt.Errorf("field %s must be a string", field)
			}
			if !pattern.MatchString(strValue) {
				return fmt.Errorf("invalid format for field: %s", field)
			}
		}

		// Sanitize fields if needed
		if sanitize, exists := r.Customer.Sanitization[field]; exists {
			strValue, ok := value.(string)
			if !ok {
				return fmt.Errorf("field %s must be a string", field)
			}
			data[field] = sanitize(strValue)
		}
	}

	// Validate age
	if dob, exists := data["date_of_birth"].(string); exists {
		birthDate, err := time.Parse("2006-01-02", dob)
		if err != nil {
			return fmt.Errorf("invalid date of birth format")
		}
		age := time.Since(birthDate).Hours() / 24 / 365.25
		if age < float64(r.Customer.MinAge) {
			return fmt.Errorf("customer must be at least %d years old", r.Customer.MinAge)
		}
	}

	return nil
}

// ValidateDocument validates document according to financial regulations
func (r *FinancialValidationRules) ValidateDocument(docType string, data []byte, uploadTime time.Time) error {
	// Check document type
	validType := false
	for _, t := range r.Document.RequiredTypes {
		if t == docType {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("invalid document type: %s", docType)
	}

	// Check file size
	if int64(len(data)) > r.Document.MaxFileSize {
		return fmt.Errorf("document size exceeds maximum allowed size of %d bytes", r.Document.MaxFileSize)
	}

	// Check document age
	if time.Since(uploadTime) > r.Document.MaxAge {
		return fmt.Errorf("document is too old, maximum age is %v", r.Document.MaxAge)
	}

	// Run content checks
	for _, check := range r.Document.ContentChecks {
		if err := check(data); err != nil {
			return fmt.Errorf("document content validation failed: %v", err)
		}
	}

	return nil
}

// ValidateTransaction validates transaction according to financial regulations
func (r *FinancialValidationRules) ValidateTransaction(data map[string]interface{}) error {
	// Check required fields
	for _, field := range r.Transaction.RequiredFields {
		if _, exists := data[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Validate amount
	amount, ok := data["amount"].(float64)
	if !ok {
		return fmt.Errorf("invalid amount format")
	}
	if amount < r.Transaction.MinAmount {
		return fmt.Errorf("amount below minimum threshold")
	}
	if amount > r.Transaction.MaxAmount {
		return fmt.Errorf("amount above maximum threshold")
	}

	// Determine risk level and run checks
	var riskLevel string
	switch {
	case amount <= r.Transaction.RiskLevels["low"].Threshold:
		riskLevel = "low"
	case amount <= r.Transaction.RiskLevels["medium"].Threshold:
		riskLevel = "medium"
	default:
		riskLevel = "high"
	}

	// Run risk level specific checks
	for _, check := range r.Transaction.RiskLevels[riskLevel].Checks {
		if err := check(data); err != nil {
			return fmt.Errorf("transaction validation failed: %v", err)
		}
	}

	return nil
}

// ValidateCompliance validates compliance requirements
func (r *FinancialValidationRules) ValidateCompliance(customerID string) error {
	// Run watchlist checks
	for _, check := range r.Compliance.WatchlistChecks {
		if err := check(customerID); err != nil {
			return fmt.Errorf("watchlist check failed: %v", err)
		}
	}

	// Run sanctions checks
	for _, check := range r.Compliance.SanctionChecks {
		if err := check(customerID); err != nil {
			return fmt.Errorf("sanctions check failed: %v", err)
		}
	}

	// Run PEP checks
	for _, check := range r.Compliance.PEPChecks {
		if err := check(customerID); err != nil {
			return fmt.Errorf("PEP check failed: %v", err)
		}
	}

	return nil
}

// CalculateRiskScore calculates overall risk score
func (r *FinancialValidationRules) CalculateRiskScore(factors map[string]float64) float64 {
	var score float64
	for factor, weight := range r.Compliance.RiskScoring {
		if value, exists := factors[factor]; exists {
			score += value * weight
		}
	}
	return score
}
