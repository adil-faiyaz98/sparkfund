package infrastructure

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TestInterfaceInjection tests for unsafe interface{} usage in JSON parsing
func TestInterfaceInjection(t *testing.T) {
	t.Run("Unsafe JSON Parsing", func(t *testing.T) {
		handler := createInterfaceHandler()

		// Test malicious JSON payloads
		payloads := []string{
			`{"data": {"__proto__": {"admin": true}}}`,
			`{"data": {"constructor": {"prototype": {"admin": true}}}}`,
			`{"data": {"toString": {"constructor": {"return": {"admin": true}}}}}`,
		}

		for _, payload := range payloads {
			req := httptest.NewRequest("POST", "/api/parse", bytes.NewBufferString(payload))
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusOK, w.Code, "Interface injection attempt should be blocked: %s", payload)
		}
	})
}

// TestGRPCSecurity tests for gRPC security vulnerabilities
func TestGRPCSecurity(t *testing.T) {
	t.Run("Unencrypted gRPC", func(t *testing.T) {
		// Test unencrypted gRPC connection
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		assert.Error(t, err, "Unencrypted gRPC connection should be rejected")
		if err == nil {
			conn.Close()
		}
	})

	t.Run("Missing mTLS", func(t *testing.T) {
		// Test missing client certificate
		creds := credentials.NewTLS(nil)
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
		assert.Error(t, err, "gRPC connection without mTLS should be rejected")
		if err == nil {
			conn.Close()
		}
	})
}

// TestReplayAttacks tests for request replay vulnerabilities
func TestReplayAttacks(t *testing.T) {
	t.Run("Token Replay", func(t *testing.T) {
		handler := createReplayHandler()

		// Create a valid request with timestamp
		req := httptest.NewRequest("POST", "/api/transaction", bytes.NewBufferString(`{"amount": 100}`))
		req.Header.Set("X-Timestamp", time.Now().Format(time.RFC3339))
		req.Header.Set("X-Nonce", "unique-nonce-1")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Valid request should succeed")

		// Replay the same request
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusOK, w.Code, "Replayed request should be rejected")
	})
}

// TestPIIExposure tests for PII/SPII exposure vulnerabilities
func TestPIIExposure(t *testing.T) {
	t.Run("Log Leakage", func(t *testing.T) {
		handler := createLogHandler()

		// Test sensitive data in logs
		req := httptest.NewRequest("GET", "/api/user/123", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		// Check if sensitive data is properly redacted in logs
		assert.NotContains(t, w.Body.String(), "credit_card", "Sensitive data should be redacted from logs")
	})

	t.Run("Response Exposure", func(t *testing.T) {
		handler := createResponseHandler()

		// Test excessive data exposure
		req := httptest.NewRequest("GET", "/api/user/profile", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotContains(t, response, "ssn", "SSN should not be exposed in response")
	})
}

// TestLogInjection tests for log injection vulnerabilities
func TestLogInjection(t *testing.T) {
	t.Run("Control Characters", func(t *testing.T) {
		handler := createLogHandler()

		// Test injection of control characters
		payloads := []string{
			"normal text\ninjected log",
			"normal text\r\ninjected log",
			"normal text\tinjected log",
			"normal text\vinjected log",
		}

		for _, payload := range payloads {
			req := httptest.NewRequest("POST", "/api/log", bytes.NewBufferString(payload))
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			assert.NotContains(t, w.Body.String(), "injected log", "Log injection should be prevented")
		}
	})
}

// Helper functions
func createInterfaceHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement safe interface{} handling
		if !isSafeInterface(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func createReplayHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement replay protection
		if !isValidReplay(r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func createLogHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement safe logging
		if !isSafeLogging(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func createResponseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement response filtering
		if !isSafeResponse(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

// Security check functions
func isSafeInterface(r *http.Request) bool {
	// Implement safe interface{} handling
	return true // Simplified for testing
}

func isValidReplay(r *http.Request) bool {
	// Implement replay protection
	return true // Simplified for testing
}

func isSafeLogging(r *http.Request) bool {
	// Implement safe logging
	return true // Simplified for testing
}

func isSafeResponse(r *http.Request) bool {
	// Implement response filtering
	return true // Simplified for testing
}
