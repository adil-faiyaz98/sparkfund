package mapper

import (
	"time"

	"github.com/google/uuid"

	"sparkfund/services/kyc-service/internal/domain"
	"sparkfund/services/kyc-service/internal/model"
)

// KYCModelToDomain converts a model.KYC to a domain.EnhancedKYC
func KYCModelToDomain(kyc *model.KYC) *domain.EnhancedKYC {
	if kyc == nil {
		return nil
	}

	var metadata domain.Metadata
	if kyc.Metadata != nil {
		metadata = domain.Metadata(kyc.Metadata)
	}

	var documentExpiry *time.Time
	if kyc.DocumentExpiry != nil {
		documentExpiry = kyc.DocumentExpiry
	}

	return &domain.EnhancedKYC{
		ID:                kyc.ID,
		UserID:            kyc.UserID,
		FirstName:         kyc.FirstName,
		LastName:          kyc.LastName,
		DateOfBirth:       kyc.DateOfBirth,
		Nationality:       kyc.Nationality,
		Email:             kyc.Email,
		PhoneNumber:       kyc.PhoneNumber,
		Address:           kyc.Address,
		City:              kyc.City,
		State:             kyc.State,
		Country:           kyc.Country,
		PostalCode:        kyc.PostalCode,
		DocumentType:      kyc.DocumentType,
		DocumentNumber:    kyc.DocumentNumber,
		DocumentFront:     kyc.DocumentFront,
		DocumentBack:      kyc.DocumentBack,
		SelfieImage:       kyc.SelfieImage,
		DocumentExpiry:    documentExpiry,
		RiskLevel:         domain.RiskLevel(kyc.RiskLevel),
		RiskScore:         kyc.RiskScore,
		TransactionAmount: kyc.TransactionAmount,
		Status:            domain.KYCStatus(kyc.Status),
		RejectionReason:   kyc.RejectionReason,
		Notes:             kyc.Notes,
		VerifiedBy:        kyc.VerifiedBy,
		VerifiedAt:        kyc.VerifiedAt,
		ReviewedBy:        kyc.ReviewedBy,
		ReviewedAt:        kyc.ReviewedAt,
		CreatedAt:         kyc.CreatedAt,
		UpdatedAt:         kyc.UpdatedAt,
		CompletedAt:       kyc.CompletedAt,
		ExpiresAt:         kyc.ExpiresAt,
		Metadata:          metadata,
	}
}

// KYCDomainToModel converts a domain.EnhancedKYC to a model.KYC
func KYCDomainToModel(kyc *domain.EnhancedKYC) *model.KYC {
	if kyc == nil {
		return nil
	}

	var metadata map[string]interface{}
	if kyc.Metadata != nil {
		metadata = map[string]interface{}(kyc.Metadata)
	}

	return &model.KYC{
		ID:                kyc.ID,
		UserID:            kyc.UserID,
		FirstName:         kyc.FirstName,
		LastName:          kyc.LastName,
		DateOfBirth:       kyc.DateOfBirth,
		Nationality:       kyc.Nationality,
		Email:             kyc.Email,
		PhoneNumber:       kyc.PhoneNumber,
		Address:           kyc.Address,
		City:              kyc.City,
		State:             kyc.State,
		Country:           kyc.Country,
		PostalCode:        kyc.PostalCode,
		DocumentType:      kyc.DocumentType,
		DocumentNumber:    kyc.DocumentNumber,
		DocumentFront:     kyc.DocumentFront,
		DocumentBack:      kyc.DocumentBack,
		SelfieImage:       kyc.SelfieImage,
		DocumentExpiry:    kyc.DocumentExpiry,
		RiskLevel:         model.RiskLevel(kyc.RiskLevel),
		RiskScore:         kyc.RiskScore,
		TransactionAmount: kyc.TransactionAmount,
		Status:            model.KYCStatus(kyc.Status),
		RejectionReason:   kyc.RejectionReason,
		Notes:             kyc.Notes,
		VerifiedBy:        kyc.VerifiedBy,
		VerifiedAt:        kyc.VerifiedAt,
		ReviewedBy:        kyc.ReviewedBy,
		ReviewedAt:        kyc.ReviewedAt,
		CreatedAt:         kyc.CreatedAt,
		UpdatedAt:         kyc.UpdatedAt,
		CompletedAt:       kyc.CompletedAt,
		ExpiresAt:         kyc.ExpiresAt,
		Metadata:          metadata,
	}
}

// KYCReviewModelToDomain converts a model.KYCReview to a domain.KYCReview
func KYCReviewModelToDomain(review *model.KYCReview) *domain.KYCReview {
	if review == nil {
		return nil
	}

	return &domain.KYCReview{
		ID:             review.ID,
		KYCID:          review.KYCID,
		ReviewerID:     review.ReviewerID,
		Status:         review.Status,
		Reason:         review.Reason,
		RiskAssessment: review.RiskAssessment,
		Notes:          review.Notes,
		CreatedAt:      review.CreatedAt,
		UpdatedAt:      review.UpdatedAt,
	}
}

// KYCReviewDomainToModel converts a domain.KYCReview to a model.KYCReview
func KYCReviewDomainToModel(review *domain.KYCReview) *model.KYCReview {
	if review == nil {
		return nil
	}

	return &model.KYCReview{
		ID:             review.ID,
		KYCID:          review.KYCID,
		ReviewerID:     review.ReviewerID,
		Status:         review.Status,
		Reason:         review.Reason,
		RiskAssessment: review.RiskAssessment,
		Notes:          review.Notes,
		CreatedAt:      review.CreatedAt,
		UpdatedAt:      review.UpdatedAt,
	}
}

// KYCModelsToDomains converts a slice of model.KYC to a slice of domain.EnhancedKYC
func KYCModelsToDomains(kycs []*model.KYC) []*domain.EnhancedKYC {
	result := make([]*domain.EnhancedKYC, len(kycs))
	for i, kyc := range kycs {
		result[i] = KYCModelToDomain(kyc)
	}
	return result
}

// KYCDomainsToModels converts a slice of domain.EnhancedKYC to a slice of model.KYC
func KYCDomainsToModels(kycs []*domain.EnhancedKYC) []*model.KYC {
	result := make([]*model.KYC, len(kycs))
	for i, kyc := range kycs {
		result[i] = KYCDomainToModel(kyc)
	}
	return result
}
