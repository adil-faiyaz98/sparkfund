package validation

import (
    "context"
    "time"

    "github.com/sparkfund/kyc-service/internal/models"
)

// DocumentValidator handles document-specific validation
type DocumentValidator struct {
    *BaseValidator
    aiClient    AIClient
    ocrService  OCRService
}

func (dv *DocumentValidator) Validate(ctx context.Context, doc *models.Document) (*ValidationResult, error) {
    startTime := time.Now()
    defer dv.metrics.RecordValidationDuration("document", time.Since(startTime))

    result := &ValidationResult{
        ValidatedAt: time.Now(),
        ModelInfo:   dv.modelInfo,
    }

    // Validate document authenticity
    authenticityScore, err := dv.validateAuthenticity(ctx, doc)
    if err != nil {
        return nil, fmt.Errorf("authenticity validation failed: %w", err)
    }

    // Validate document content
    contentScore, err := dv.validateContent(ctx, doc)
    if err != nil {
        return nil, fmt.Errorf("content validation failed: %w", err)
    }

    // Calculate final score
    result.Score = (authenticityScore + contentScore) / 2
    result.IsValid = result.Score >= dv.config.MinDocumentScore

    return result, nil
}