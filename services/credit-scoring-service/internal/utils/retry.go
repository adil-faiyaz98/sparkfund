package utils

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/sparkfund/credit-scoring-service/internal/errors"
	"go.uber.org/zap"
)

// RetryConfig defines the configuration for retry operations
type RetryConfig struct {
	MaxAttempts      int
	InitialDelay     time.Duration
	MaxDelay         time.Duration
	BackoffFactor    float64
	JitterFactor     float64
	RetryableErrors  []error
	NonRetryableErrors []error
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:      3,
		InitialDelay:     time.Second,
		MaxDelay:         30 * time.Second,
		BackoffFactor:    2.0,
		JitterFactor:     0.1,
		RetryableErrors: []error{
			errors.NewAPIError(errors.ErrExternalService, ""),
			errors.NewAPIError(errors.ErrDatabase, ""),
		},
		NonRetryableErrors: []error{
			errors.NewAPIError(errors.ErrValidation, ""),
			errors.NewAPIError(errors.ErrUnauthorized, ""),
			errors.NewAPIError(errors.ErrForbidden, ""),
		},
	}
}

// Retry executes a function with retry logic and exponential backoff
func Retry(ctx context.Context, config *RetryConfig, logger *zap.Logger, operation func() error) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// Check context cancellation
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled: %w", err)
		}

		// Execute operation
		err := operation()
		if err == nil {
			return nil
		}
		lastErr = err

		// Check if error is non-retryable
		if isNonRetryableError(err, config.NonRetryableErrors) {
			return err
		}

		// Check if error is retryable
		if !isRetryableError(err, config.RetryableErrors) {
			return err
		}

		// Calculate next delay with exponential backoff and jitter
		delay = calculateNextDelay(delay, config)

		// Log retry attempt
		logger.Warn("retry attempt",
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", config.MaxAttempts),
			zap.Duration("delay", delay),
			zap.Error(err),
		)

		// Wait before next attempt
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		case <-time.After(delay):
			continue
		}
	}

	return fmt.Errorf("max retry attempts reached: %w", lastErr)
}

// isRetryableError checks if an error is retryable based on the configuration
func isRetryableError(err error, retryableErrors []error) bool {
	for _, retryableErr := range retryableErrors {
		if isErrorType(err, retryableErr) {
			return true
		}
	}
	return false
}

// isNonRetryableError checks if an error is non-retryable based on the configuration
func isNonRetryableError(err error, nonRetryableErrors []error) bool {
	for _, nonRetryableErr := range nonRetryableErrors {
		if isErrorType(err, nonRetryableErr) {
			return true
		}
	}
	return false
}

// isErrorType checks if an error is of a specific type
func isErrorType(err, targetErr error) bool {
	if apiErr, ok := err.(*errors.APIError); ok {
		if targetAPIErr, ok := targetErr.(*errors.APIError); ok {
			return apiErr.Code == targetAPIErr.Code
		}
	}
	return false
}

// calculateNextDelay calculates the next delay with exponential backoff and jitter
func calculateNextDelay(currentDelay time.Duration, config *RetryConfig) time.Duration {
	// Apply exponential backoff
	nextDelay := float64(currentDelay) * config.BackoffFactor

	// Apply jitter
	jitter := nextDelay * config.JitterFactor
	nextDelay += (math.Rand.Float64()*2 - 1) * jitter

	// Ensure delay doesn't exceed max delay
	if nextDelay > float64(config.MaxDelay) {
		nextDelay = float64(config.MaxDelay)
	}

	return time.Duration(nextDelay)
} 