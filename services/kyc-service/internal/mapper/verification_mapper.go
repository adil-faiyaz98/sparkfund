package mapper

import (
	"sparkfund/services/kyc-service/internal/domain"
	"sparkfund/services/kyc-service/internal/model"
)

// VerificationModelToDomain converts a model.Verification to a domain.EnhancedVerification
func VerificationModelToDomain(ver *model.Verification) *domain.EnhancedVerification {
	if ver == nil {
		return nil
	}

	var metadata domain.Metadata
	if ver.Metadata != nil {
		metadata = domain.Metadata(ver.Metadata)
	}

	var result domain.Metadata
	if ver.Result != nil {
		result = domain.Metadata(ver.Result)
	}

	return &domain.EnhancedVerification{
		ID:              ver.ID,
		KYCID:           ver.KYCID,
		DocumentID:      ver.DocumentID,
		Type:            domain.VerificationType(ver.Type),
		Status:          domain.VerificationStatus(ver.Status),
		Method:          domain.VerificationMethod(ver.Method),
		VerifierID:      ver.VerifierID,
		ConfidenceScore: ver.ConfidenceScore,
		MatchScore:      ver.MatchScore,
		FraudScore:      ver.FraudScore,
		Notes:           ver.Notes,
		Metadata:        metadata,
		Result:          result,
		ErrorMessage:    ver.ErrorMessage,
		CreatedAt:       ver.CreatedAt,
		UpdatedAt:       ver.UpdatedAt,
		CompletedAt:     ver.CompletedAt,
		ExpiresAt:       ver.ExpiresAt,
	}
}

// VerificationDomainToModel converts a domain.EnhancedVerification to a model.Verification
func VerificationDomainToModel(ver *domain.EnhancedVerification) *model.Verification {
	if ver == nil {
		return nil
	}

	var metadata map[string]interface{}
	if ver.Metadata != nil {
		metadata = map[string]interface{}(ver.Metadata)
	}

	var result map[string]interface{}
	if ver.Result != nil {
		result = map[string]interface{}(ver.Result)
	}

	return &model.Verification{
		ID:              ver.ID,
		KYCID:           ver.KYCID,
		DocumentID:      ver.DocumentID,
		Type:            model.VerificationType(ver.Type),
		Status:          model.VerificationStatus(ver.Status),
		Method:          model.VerificationMethod(ver.Method),
		VerifierID:      ver.VerifierID,
		ConfidenceScore: ver.ConfidenceScore,
		MatchScore:      ver.MatchScore,
		FraudScore:      ver.FraudScore,
		Notes:           ver.Notes,
		Metadata:        metadata,
		Result:          result,
		ErrorMessage:    ver.ErrorMessage,
		CreatedAt:       ver.CreatedAt,
		UpdatedAt:       ver.UpdatedAt,
		CompletedAt:     ver.CompletedAt,
		ExpiresAt:       ver.ExpiresAt,
	}
}

// VerificationResultModelToDomain converts a model.VerificationResult to a domain.VerificationResult
func VerificationResultModelToDomain(result *model.VerificationResult) *domain.VerificationResult {
	if result == nil {
		return nil
	}

	var details domain.Metadata
	if result.Details != nil {
		details = domain.Metadata(result.Details)
	}

	return &domain.VerificationResult{
		ID:             result.ID,
		VerificationID: result.VerificationID,
		Score:          result.Score,
		MatchScore:     result.MatchScore,
		FraudScore:     result.FraudScore,
		Success:        result.Success,
		Details:        details,
		ErrorMessage:   result.ErrorMessage,
		CreatedAt:      result.CreatedAt,
	}
}

// VerificationResultDomainToModel converts a domain.VerificationResult to a model.VerificationResult
func VerificationResultDomainToModel(result *domain.VerificationResult) *model.VerificationResult {
	if result == nil {
		return nil
	}

	var details map[string]interface{}
	if result.Details != nil {
		details = map[string]interface{}(result.Details)
	}

	return &model.VerificationResult{
		ID:             result.ID,
		VerificationID: result.VerificationID,
		Score:          result.Score,
		MatchScore:     result.MatchScore,
		FraudScore:     result.FraudScore,
		Success:        result.Success,
		Details:        details,
		ErrorMessage:   result.ErrorMessage,
		CreatedAt:      result.CreatedAt,
		UpdatedAt:      result.CreatedAt, // Use CreatedAt as UpdatedAt for new records
	}
}

// VerificationSummaryModelToDomain converts a model.VerificationSummary to a domain.VerificationSummary
func VerificationSummaryModelToDomain(summary *model.VerificationSummary) *domain.VerificationSummary {
	if summary == nil {
		return nil
	}

	return &domain.VerificationSummary{
		ID:              summary.ID,
		KYCID:           summary.KYCID,
		DocumentID:      summary.DocumentID,
		Type:            domain.VerificationType(summary.Type),
		Status:          domain.VerificationStatus(summary.Status),
		Method:          domain.VerificationMethod(summary.Method),
		ConfidenceScore: summary.ConfidenceScore,
		MatchScore:      summary.MatchScore,
		FraudScore:      summary.FraudScore,
		CreatedAt:       summary.CreatedAt,
		CompletedAt:     summary.CompletedAt,
		ProcessingTime:  summary.ProcessingTime,
		Success:         summary.Success,
	}
}

// VerificationModelsToDomains converts a slice of model.Verification to a slice of domain.EnhancedVerification
func VerificationModelsToDomains(verifications []*model.Verification) []*domain.EnhancedVerification {
	result := make([]*domain.EnhancedVerification, len(verifications))
	for i, ver := range verifications {
		result[i] = VerificationModelToDomain(ver)
	}
	return result
}

// VerificationDomainsToModels converts a slice of domain.EnhancedVerification to a slice of model.Verification
func VerificationDomainsToModels(verifications []*domain.EnhancedVerification) []*model.Verification {
	result := make([]*model.Verification, len(verifications))
	for i, ver := range verifications {
		result[i] = VerificationDomainToModel(ver)
	}
	return result
}
