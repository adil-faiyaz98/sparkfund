package models

// ErrorResponse represents a standardized error response format for API errors
// @Description Error response
type ErrorResponse struct {
	Error string `json:"error" example:"Error message description"`
}

// SuccessResponse represents a standardized success response with a message
// @Description Success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}
