package dto

import (
	"time"

	"github.com/google/uuid"
	"sparkfund/services/kyc-service/internal/domain"
)

// VerificationRequest represents a request to create a verification
type VerificationRequest struct {
	DocumentID uuid.UUID `json:"document_id" binding:"required"`
	Method     string    `json:"method" binding:"required"`
}

// VerificationResponse represents a verification response
type VerificationResponse struct {
	ID              uuid.UUID `json:"id"`
	KYCID           uuid.UUID `json:"kyc_id,omitempty"`
	DocumentID      uuid.UUID `json:"document_id,omitempty"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Method          string    `json:"method"`
	ConfidenceScore float64   `json:"confidence_score"`
	MatchScore      float64   `json:"match_score,omitempty"`
	FraudScore      float64   `json:"fraud_score,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CompletedAt     string    `json:"completed_at,omitempty"`
}

// VerificationListResponse represents a paginated list of verifications
type VerificationListResponse struct {
	Verifications []VerificationResponse `json:"verifications"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
}

// VerificationStatusUpdateRequest represents a request to update verification status
type VerificationStatusUpdateRequest struct {
	Status          string  `json:"status" binding:"required"`
	ConfidenceScore float64 `json:"confidence_score"`
	Notes           string  `json:"notes,omitempty"`
}

// VerificationResultRequest represents a request to create a verification result
type VerificationResultRequest struct {
	Score        float64                 `json:"score" binding:"required"`
	MatchScore   float64                 `json:"match_score,omitempty"`
	FraudScore   float64                 `json:"fraud_score,omitempty"`
	Success      bool                    `json:"success" binding:"required"`
	Details      map[string]interface{}  `json:"details,omitempty"`
	ErrorMessage string                  `json:"error_message,omitempty"`
}

// FromDomainVerification converts a domain verification to a verification response
func FromDomainVerification(ver *domain.EnhancedVerification) VerificationResponse {
	response := VerificationResponse{
		ID:              ver.ID,
		Type:            string(ver.Type),
		Status:          string(ver.Status),
		Method:          string(ver.Method),
		ConfidenceScore: ver.ConfidenceScore,
		MatchScore:      ver.MatchScore,
		FraudScore:      ver.FraudScore,
		Notes:           ver.Notes,
		CreatedAt:       ver.CreatedAt,
		UpdatedAt:       ver.UpdatedAt,
	}

	if ver.KYCID != nil {
		response.KYCID = *ver.KYCID
	}

	if ver.DocumentID != nil {
		response.DocumentID = *ver.DocumentID
	}

	if ver.CompletedAt != nil {
		response.CompletedAt = ver.CompletedAt.Format("2006-01-02T15:04:05Z")
	}

	return response
}

// FromDomainVerifications converts a slice of domain verifications to verification responses
func FromDomainVerifications(verifications []*domain.EnhancedVerification) []VerificationResponse {
	responses := make([]VerificationResponse, len(verifications))
	for i, ver := range verifications {
		responses[i] = FromDomainVerification(ver)
	}
	return responses
}

// ToDomainVerificationResult converts a verification result request to a domain verification result
func ToDomainVerificationResult(req *VerificationResultRequest) *domain.VerificationResult {
	return &domain.VerificationResult{
		ID:           uuid.New(),
		Score:        req.Score,
		MatchScore:   req.MatchScore,
		FraudScore:   req.FraudScore,
		Success:      req.Success,
		Details:      req.Details,
		ErrorMessage: req.ErrorMessage,
		CreatedAt:    time.Now(),
	}
}
