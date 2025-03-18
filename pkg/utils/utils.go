package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewResponse creates a new success response
func NewResponse(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err error) Response {
	return Response{
		Success: false,
		Error:   err.Error(),
	}
}

// WriteJSON writes a JSON response to the http.ResponseWriter
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, int(math.Ceil(float64(length)*3/4)))
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// FormatMoney formats a float64 as a money string with 2 decimal places
func FormatMoney(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

// ParseMoney parses a money string into a float64
func ParseMoney(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "$")
	var amount float64
	_, err := fmt.Sscanf(s, "%f", &amount)
	return amount, err
}

// TimeAgo returns a human-readable string representing time since the given time
func TimeAgo(t time.Time) string {
	duration := time.Since(t)
	hours := duration.Hours()

	switch {
	case hours < 1:
		minutes := int(duration.Minutes())
		if minutes < 1 {
			return "just now"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case hours < 24:
		h := int(hours)
		return fmt.Sprintf("%d hours ago", h)
	case hours < 168: // 7 days
		days := int(hours / 24)
		return fmt.Sprintf("%d days ago", days)
	case hours < 720: // 30 days
		weeks := int(hours / 168)
		return fmt.Sprintf("%d weeks ago", weeks)
	case hours < 8760: // 365 days
		months := int(hours / 720)
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(hours / 8760)
		return fmt.Sprintf("%d years ago", years)
	}
}

// Truncate truncates a string to the specified length and adds ellipsis if needed
func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

// Contains checks if a string slice contains a specific string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Map applies a function to each element in a slice and returns a new slice
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// Filter returns a new slice containing only the elements that satisfy the predicate
func Filter[T any](slice []T, f func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce reduces a slice to a single value using an accumulator function
func Reduce[T, U any](slice []T, initial U, f func(U, T) U) U {
	result := initial
	for _, v := range slice {
		result = f(result, v)
	}
	return result
}

// Example usage:
// response := utils.NewResponse(data)
// utils.WriteJSON(w, http.StatusOK, response)
//
// randomString, _ := utils.GenerateRandomString(32)
// formattedAmount := utils.FormatMoney(123.456)
// timeAgo := utils.TimeAgo(someTime)
//
// numbers := []int{1, 2, 3, 4, 5}
// doubled := utils.Map(numbers, func(n int) int { return n * 2 })
// evens := utils.Filter(numbers, func(n int) bool { return n%2 == 0 })
// sum := utils.Reduce(numbers, 0, func(acc, n int) int { return acc + n })
