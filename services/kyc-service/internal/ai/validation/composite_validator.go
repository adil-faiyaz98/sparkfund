package validation

import (
    "context"
    "sync"

    "github.com/sparkfund/kyc-service/internal/models"
)

// CompositeValidator orchestrates multiple validators
type CompositeValidator struct {
    *BaseValidator
    documentValidator  *DocumentValidator
    biometricValidator *BiometricValidator
    threatValidator    *ThreatValidator
    behaviorValidator  *BehaviorValidator
}

// NewCompositeValidator creates a new composite validator
func NewCompositeValidator(config *Config) *CompositeValidator {
    return &CompositeValidator{
        BaseValidator:      NewBaseValidator(config),
        documentValidator:  NewDocumentValidator(config),
        biometricValidator: NewBiometricValidator(config),
        threatValidator:    NewThreatValidator(config),
        behaviorValidator:  NewBehaviorValidator(config),
    }
}

// ValidateKYC performs comprehensive KYC validation
func (cv *CompositeValidator) ValidateKYC(ctx context.Context, input *models.KYCData) (*ValidationResult, error) {
    var wg sync.WaitGroup
    results := make(chan *ValidationResult, 4)
    errors := make(chan error, 4)

    // Run all validations in parallel
    validators := []struct {
        name string
        fn   func() (*ValidationResult, error)
    }{
        {"document", func() (*ValidationResult, error) { 
            return cv.documentValidator.Validate(ctx, input.Document) 
        }},
        {"biometric", func() (*ValidationResult, error) { 
            return cv.biometricValidator.Validate(ctx, input.BiometricData) 
        }},
        {"threat", func() (*ValidationResult, error) { 
            return cv.threatValidator.Validate(ctx, input) 
        }},
        {"behavior", func() (*ValidationResult, error) { 
            return cv.behaviorValidator.Validate(ctx, input.BehaviorData) 
        }},
    }

    for _, v := range validators {
        wg.Add(1)
        go func(name string, validatorFn func() (*ValidationResult, error)) {
            defer wg.Done()
            result, err := validatorFn()
            if err != nil {
                errors <- err
                return
            }
            results <- result
        }(v.name, v.fn)
    }

    // Wait for all validations to complete
    go func() {
        wg.Wait()
        close(results)
        close(errors)
    }()

    // Aggregate results
    return cv.aggregateResults(results, errors)
}

// aggregateResults combines multiple validation results
func (cv *CompositeValidator) aggregateResults(results <-chan *ValidationResult, errors <-chan error) (*ValidationResult, error) {
    var finalResult ValidationResult
    var errs []error

    for err := range errors {
        errs = append(errs, err)
    }

    if len(errs) > 0 {
        return nil, fmt.Errorf("validation errors occurred: %v", errs)
    }

    var totalScore float64
    var count int
    for result := range results {
        totalScore += result.Score
        count++
        finalResult.Warnings = append(finalResult.Warnings, result.Warnings...)
    }

    finalResult.Score = totalScore / float64(count)
    finalResult.IsValid = finalResult.Score >= cv.config.MinValidationScore
    finalResult.ValidatedAt = time.Now()

    return &finalResult, nil
}