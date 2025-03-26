package models

// Error represents a standardized error response
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}