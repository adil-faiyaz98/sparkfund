package models

import "time"

// TokenValidationRequest holds data for token validation requests
type TokenValidationRequest struct {
	Token string `json:"token" validate:"required"`
}

// TokenValidationResponse is the response from token validation
type TokenValidationResponse struct {
	Valid     bool                   `json:"valid"`
	Subject   string                 `json:"subject,omitempty"`
	Issuer    string                 `json:"issuer,omitempty"`
	IssuedAt  time.Time              `json:"issued_at,omitempty"`
	ExpiresAt time.Time              `json:"expires_at,omitempty"`
	Claims    map[string]interface{} `json:"claims,omitempty"`
}

// TokenGenerationRequest is the request for token generation
type TokenGenerationRequest struct {
	Subject   string                 `json:"subject" validate:"required"`
	ExpiresIn int64                  `json:"expires_in,omitempty"` // Duration in seconds
	Claims    map[string]interface{} `json:"claims,omitempty"`
}

// TokenGenerationResponse is the response for token generation
type TokenGenerationResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// TokenRefreshRequest is the request for refreshing tokens
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenRefreshResponse is the response for token refresh
type TokenRefreshResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
}
