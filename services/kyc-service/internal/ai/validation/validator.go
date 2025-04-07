package validation

import (
    "context"
    "errors"
    "time"

    "github.com/sparkfund/kyc-service/internal/models"
)

// ValidationResult represents the outcome of a validation check
type ValidationResult struct {
    IsValid     bool
    Score       float64
    Confidence  float64
    Errors      []error
    Warnings    []string
    ValidatedAt time.Time
    ModelInfo   ModelInfo
}

// ModelInfo contains information about the model used for validation
type ModelInfo struct {
    Name        string
    Version     string
    LastUpdated time.Time
}

// Validator interface defines the contract for all validators
type Validator interface {
    Validate(ctx context.Context, input interface{}) (*ValidationResult, error)
    ValidateAsync(ctx context.Context, input interface{}) <-chan *ValidationResult
    GetConfidence() float64
}

// BaseValidator provides common validation functionality
type BaseValidator struct {
    config     *Config
    metrics    *Metrics
    modelInfo  ModelInfo
    confidence float64
}

// NewBaseValidator creates a new base validator
func NewBaseValidator(config *Config) *BaseValidator {
    return &BaseValidator{
        config:  config,
        metrics: NewMetrics(),
    }
}

// ValidateWithRetry performs validation with retry logic
func (v *BaseValidator) ValidateWithRetry(ctx context.Context, input interface{}, maxRetries int) (*ValidationResult, error) {
    var lastErr error
    for i := 0; i < maxRetries; i++ {
        result, err := v.Validate(ctx, input)
        if err == nil {
            return result, nil
        }
        if errors.Is(err, context.Canceled) {
            return nil, err
        }
        lastErr = err
        time.Sleep(v.config.RetryDelay)
    }
    return nil, fmt.Errorf("validation failed after %d retries: %w", maxRetries, lastErr)
}

// ValidateAsync performs asynchronous validation
func (v *BaseValidator) ValidateAsync(ctx context.Context, input interface{}) <-chan *ValidationResult {
    resultChan := make(chan *ValidationResult, 1)
    go func() {
        defer close(resultChan)
        result, err := v.Validate(ctx, input)
        if err != nil {
            result = &ValidationResult{
                IsValid:     false,
                Errors:      []error{err},
                ValidatedAt: time.Now(),
            }
        }
        select {
        case <-ctx.Done():
            return
        case resultChan <- result:
        }
    }()
    return resultChan
}