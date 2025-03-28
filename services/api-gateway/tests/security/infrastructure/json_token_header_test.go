package infrastructure

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

// TestUnsafeJSONParsing tests for unsafe interface{} usage in JSON parsing
func TestUnsafeJSONParsing(t *testing.T) {
	t.Run("Type Confusion Attacks", func(t *testing.T) {
		payloads := []map[string]interface{}{
			{"isAdmin": "true"}, // String instead of bool
			{"role": 1},         // Number instead of string
			{"permissions": []interface{}{"admin", 1, true}}, // Mixed types
			{"enabled": "1"},                        // String "1" instead of bool
			{"level": "5"},                          // String number instead of int
			{"tags": []interface{}{1, "tag", true}}, // Mixed array types
		}

		for _, payload := range payloads {
			jsonData, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", "/api/auth", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler := createJSONHandler()
			handler.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusOK, w.Code, "Type confusion attack should be blocked: %v", payload)
		}
	})

	t.Run("Interface{} Abuse", func(t *testing.T) {
		// Test unmarshaling into map[string]interface{}
		var result map[string]interface{}
		payload := `{"data": {"nested": {"key": "value"}}}`
		err := json.Unmarshal([]byte(payload), &result)
		assert.Error(t, err, "Should not allow unmarshaling into interface{}")

		// Test with strict struct
		var strictResult struct {
			Data struct {
				Nested struct {
					Key string `json:"key"`
				} `json:"nested"`
			} `json:"data"`
		}
		err = json.Unmarshal([]byte(payload), &strictResult)
		assert.NoError(t, err, "Should allow unmarshaling into strict struct")
	})
}

// TestJWTTokenForwarding tests for unverified JWT token forwarding
func TestJWTTokenForwarding(t *testing.T) {
	t.Run("Token Validation Bypass", func(t *testing.T) {
		// Create a forged token without proper validation
		token := createForgedToken()

		// Test token forwarding without validation
		req := httptest.NewRequest("GET", "/api/internal/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		handler := createJWTHandler()
		handler.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code, "Unverified token should be rejected")
	})

	t.Run("Token Forgery", func(t *testing.T) {
		// Test with various forged token scenarios
		tokens := []string{
			"eyJhbGciOiJub25lIn0.eyJpc3MiOiJhZG1pbiJ9.",                                                                       // No signature
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhZG1pbiJ9.invalid",                                               // Invalid signature
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhZG1pbiIsImlhdCI6MTYxNjI0ODAwMCwiZXhwIjoxNjE2MjQ4MDAwfQ.invalid", // Expired token
		}

		for _, token := range tokens {
			req := httptest.NewRequest("GET", "/api/internal/users", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			handler := createJWTHandler()
			handler.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusOK, w.Code, "Forged token should be rejected: %s", token)
		}
	})
}

// TestHeaderTrust tests for improper header trust between services
func TestHeaderTrust(t *testing.T) {
	t.Run("Header Injection", func(t *testing.T) {
		headers := map[string]string{
			"X-User-ID":       "admin",
			"X-Role":          "superuser",
			"X-Forwarded-For": "internal-service",
			"X-Org-ID":        "123",
			"X-Auth-Token":    "fake-token",
		}

		for key, value := range headers {
			req := httptest.NewRequest("GET", "/api/admin", nil)
			req.Header.Set(key, value)
			w := httptest.NewRecorder()
			handler := createHeaderHandler()
			handler.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusOK, w.Code, "Header injection attempt should be blocked: %s", key)
		}
	})

	t.Run("Gateway Bypass", func(t *testing.T) {
		// Test direct access to backend service
		req := httptest.NewRequest("GET", "/api/internal/users", nil)
		req.Header.Set("X-User-ID", "admin")
		req.Header.Set("X-Role", "superuser")
		w := httptest.NewRecorder()
		handler := createHeaderHandler()
		handler.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code, "Direct backend access should be blocked")
	})
}

// Helper functions
func createJSONHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement strict JSON parsing
		var strictStruct struct {
			IsAdmin     bool     `json:"isAdmin"`
			Role        string   `json:"role"`
			Permissions []string `json:"permissions"`
		}

		if err := json.NewDecoder(r.Body).Decode(&strictStruct); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func createJWTHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement proper JWT validation
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your-secret-key"), nil // Use proper key management
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func createHeaderHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement header validation
		if !isValidGatewayRequest(r) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func createForgedToken() string {
	// Create a forged token for testing
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "admin"
	claims["iat"] = 1616248000
	claims["exp"] = 1616248000
	tokenString, _ := token.SignedString([]byte("wrong-secret"))
	return tokenString
}

func isValidGatewayRequest(r *http.Request) bool {
	// Implement gateway validation
	return r.Header.Get("X-Gateway-ID") == "trusted-gateway"
}
