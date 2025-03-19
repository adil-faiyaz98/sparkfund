package errors

import "errors"

var (
	ErrNoInstancesAvailable = errors.New("no service instances available")
	ErrInvalidToken         = errors.New("invalid token")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrServiceUnavailable   = errors.New("service unavailable")
	ErrInternalServer       = errors.New("internal server error")
)
