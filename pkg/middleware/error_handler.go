package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yourusername/money-pulse/pkg/errors"
)

// ErrorHandler middleware for consistent error responses
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture status
		crw := &customResponseWriter{ResponseWriter: w}

		// Process the request
		next.ServeHTTP(crw, r)
	})
}

type customResponseWriter struct {
	http.ResponseWriter
	status int
}

func (crw *customResponseWriter) WriteHeader(status int) {
	crw.status = status
	crw.ResponseWriter.WriteHeader(status)
}

// RespondWithError sends a consistent error response
func RespondWithError(w http.ResponseWriter, err error) {
	var appErr *errors.AppError

	// Check if it's already an AppError
	if errors.As(err, &appErr) {
		sendJSONResponse(w, appErr.Code, appErr)
		return
	}

	// Default to internal server error
	appErr = errors.NewInternalError(err)
	sendJSONResponse(w, appErr.Code, appErr)
}

func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}
