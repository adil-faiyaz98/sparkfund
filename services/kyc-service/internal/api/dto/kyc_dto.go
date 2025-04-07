package dto

import (
	"time"

	"github.com/google/uuid"
	"sparkfund/services/kyc-service/internal/domain"
)

// KYCRequest represents a request to create a KYC verification
type KYCRequest struct {
	FirstName         string  `json:"first_name" binding:"required"`
	LastName          string  `json:"last_name" binding:"required"`
	DateOfBirth       string  `json:"date_of_birth" binding:"required"` // Format: YYYY-MM-DD
	Nationality       string  `json:"nationality,omitempty"`
	Email             string  `json:"email,omitempty"`
	PhoneNumber       string  `json:"phone_number,omitempty"`
	Address           string  `json:"address" binding:"required"`
	City              string  `json:"city" binding:"required"`
	State             string  `json:"state,omitempty"`
	Country           string  `json:"country" binding:"required"`
	PostalCode        string  `json:"postal_code" binding:"required"`
	DocumentType      string  `json:"document_type" binding:"required"`
	DocumentNumber    string  `json:"document_number" binding:"required"`
	DocumentFront     string  `json:"document_front" binding:"required"`
	DocumentBack      string  `json:"document_back" binding:"required"`
	SelfieImage       string  `json:"selfie_image" binding:"required"`
	DocumentExpiry    string  `json:"document_expiry,omitempty"` // Format: YYYY-MM-DD
	TransactionAmount float64 `json:"transaction_amount,omitempty"`
}

// KYCResponse represents a response for a KYC verification
type KYCResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	Status          string    `json:"status"`
	RiskLevel       string    `json:"risk_level"`
	RiskScore       float64   `json:"risk_score"`
	RejectionReason string    `json:"rejection_reason,omitempty"`
	VerifiedAt      string    `json:"verified_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CompletedAt     string    `json:"completed_at,omitempty"`
}

// KYCDetailResponse represents a detailed response for a KYC verification
type KYCDetailResponse struct {
	ID              uuid.UUID           `json:"id"`
	UserID          uuid.UUID           `json:"user_id"`
	FirstName       string              `json:"first_name"`
	LastName        string              `json:"last_name"`
	DateOfBirth     string              `json:"date_of_birth"`
	Nationality     string              `json:"nationality,omitempty"`
	Email           string              `json:"email,omitempty"`
	PhoneNumber     string              `json:"phone_number,omitempty"`
	Address         string              `json:"address"`
	City            string              `json:"city"`
	State           string              `json:"state,omitempty"`
	Country         string              `json:"country"`
	PostalCode      string              `json:"postal_code"`
	Status          string              `json:"status"`
	RiskLevel       string              `json:"risk_level"`
	RiskScore       float64             `json:"risk_score"`
	RejectionReason string              `json:"rejection_reason,omitempty"`
	Notes           string              `json:"notes,omitempty"`
	Documents       []DocumentResponse  `json:"documents,omitempty"`
	Verifications   []VerificationResponse `json:"verifications,omitempty"`
	VerifiedAt      string              `json:"verified_at,omitempty"`
	ReviewedAt      string              `json:"reviewed_at,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	CompletedAt     string              `json:"completed_at,omitempty"`
}

// KYCListResponse represents a paginated list of KYC verifications
type KYCListResponse struct {
	KYCs     []KYCResponse `json:"kycs"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// KYCStatusUpdateRequest represents a request to update KYC status
type KYCStatusUpdateRequest struct {
	Status     string    `json:"status" binding:"required"`
	Notes      string    `json:"notes,omitempty"`
	ReviewerID uuid.UUID `json:"reviewer_id" binding:"required"`
}

// KYCRiskUpdateRequest represents a request to update KYC risk level
type KYCRiskUpdateRequest struct {
	RiskLevel  string    `json:"risk_level" binding:"required"`
	RiskScore  float64   `json:"risk_score" binding:"required"`
	Notes      string    `json:"notes,omitempty"`
	ReviewerID uuid.UUID `json:"reviewer_id" binding:"required"`
}

// FromDomainKYC converts a domain KYC to a KYC response
func FromDomainKYC(kyc *domain.EnhancedKYC) KYCResponse {
	response := KYCResponse{
		ID:              kyc.ID,
		UserID:          kyc.UserID,
		Status:          string(kyc.Status),
		RiskLevel:       string(kyc.RiskLevel),
		RiskScore:       kyc.RiskScore,
		RejectionReason: kyc.RejectionReason,
		CreatedAt:       kyc.CreatedAt,
		UpdatedAt:       kyc.UpdatedAt,
	}

	if kyc.VerifiedAt != nil {
		response.VerifiedAt = kyc.VerifiedAt.Format("2006-01-02T15:04:05Z")
	}

	if kyc.CompletedAt != nil {
		response.CompletedAt = kyc.CompletedAt.Format("2006-01-02T15:04:05Z")
	}

	return response
}

// FromDomainKYCDetail converts a domain KYC to a detailed KYC response
func FromDomainKYCDetail(kyc *domain.EnhancedKYC) KYCDetailResponse {
	response := KYCDetailResponse{
		ID:              kyc.ID,
		UserID:          kyc.UserID,
		FirstName:       kyc.FirstName,
		LastName:        kyc.LastName,
		DateOfBirth:     kyc.DateOfBirth,
		Nationality:     kyc.Nationality,
		Email:           kyc.Email,
		PhoneNumber:     kyc.PhoneNumber,
		Address:         kyc.Address,
		City:            kyc.City,
		State:           kyc.State,
		Country:         kyc.Country,
		PostalCode:      kyc.PostalCode,
		Status:          string(kyc.Status),
		RiskLevel:       string(kyc.RiskLevel),
		RiskScore:       kyc.RiskScore,
		RejectionReason: kyc.RejectionReason,
		Notes:           kyc.Notes,
		CreatedAt:       kyc.CreatedAt,
		UpdatedAt:       kyc.UpdatedAt,
	}

	if kyc.VerifiedAt != nil {
		response.VerifiedAt = kyc.VerifiedAt.Format("2006-01-02T15:04:05Z")
	}

	if kyc.ReviewedAt != nil {
		response.ReviewedAt = kyc.ReviewedAt.Format("2006-01-02T15:04:05Z")
	}

	if kyc.CompletedAt != nil {
		response.CompletedAt = kyc.CompletedAt.Format("2006-01-02T15:04:05Z")
	}

	if len(kyc.Documents) > 0 {
		response.Documents = FromDomainDocuments(kyc.Documents)
	}

	if len(kyc.Verifications) > 0 {
		response.Verifications = FromDomainVerifications(kyc.Verifications)
	}

	return response
}

// FromDomainKYCs converts a slice of domain KYCs to KYC responses
func FromDomainKYCs(kycs []*domain.EnhancedKYC) []KYCResponse {
	responses := make([]KYCResponse, len(kycs))
	for i, kyc := range kycs {
		responses[i] = FromDomainKYC(kyc)
	}
	return responses
}
