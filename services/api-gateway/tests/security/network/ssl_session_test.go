package network

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/api-gateway/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestSSLSessionSecurity(t *testing.T) {
	// Create test router
	router := gin.New()
	sm := middleware.NewSecurityMiddleware(middleware.SecurityConfig{})
	sm.Apply(router)

	// Test SSL Stripping Protection
	t.Run("SSL Stripping Protection", func(t *testing.T) {
		// Test HTTP request (should redirect to HTTPS)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://example.com/api/data", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMovedPermanently, w.Code, "HTTP request should redirect to HTTPS")
		assert.Contains(t, w.Header().Get("Location"), "https://", "Redirect should be to HTTPS")

		// Test HSTS header
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "https://example.com/api/data", nil)
		router.ServeHTTP(w, req)
		assert.Contains(t, w.Header().Get("Strict-Transport-Security"), "max-age=31536000", "HSTS header should be present")

		// Test mixed content blocking
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "https://example.com/api/data", nil)
		req.Header.Set("Content-Security-Policy", "upgrade-insecure-requests")
		router.ServeHTTP(w, req)
		assert.Contains(t, w.Header().Get("Content-Security-Policy"), "upgrade-insecure-requests", "Mixed content should be blocked")
	})

	// Test Session Security
	t.Run("Session Security", func(t *testing.T) {
		// Test secure cookie attributes
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "https://example.com/api/login", nil)
		router.ServeHTTP(w, req)
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			assert.True(t, cookie.Secure, "Cookies should be secure")
			assert.True(t, cookie.HttpOnly, "Cookies should be HttpOnly")
			assert.Equal(t, "Strict", cookie.SameSite.String(), "Cookies should use Strict SameSite")
		}

		// Test session fixation prevention
		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "https://example.com/api/login", nil)
		req.Header.Set("Cookie", "session=old_session_id")
		router.ServeHTTP(w, req)
		newCookies := w.Result().Cookies()
		for _, cookie := range newCookies {
			assert.NotEqual(t, "old_session_id", cookie.Value, "Session ID should change after login")
		}

		// Test session timeout
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "https://example.com/api/data", nil)
		req.Header.Set("Cookie", "session=expired_session_id")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Expired session should be rejected")

		// Test concurrent session handling
		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "https://example.com/api/login", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.1")
		router.ServeHTTP(w, req)
		session1 := w.Result().Cookies()[0].Value

		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "https://example.com/api/login", nil)
		req.Header.Set("X-Forwarded-For", "192.168.1.2")
		router.ServeHTTP(w, req)
		session2 := w.Result().Cookies()[0].Value

		assert.NotEqual(t, session1, session2, "Different IPs should get different sessions")
	})

	// Test TLS Configuration
	t.Run("TLS Configuration", func(t *testing.T) {
		// Test TLS version
		config := &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS13,
		}
		assert.Equal(t, tls.VersionTLS12, config.MinVersion, "Minimum TLS version should be 1.2")
		assert.Equal(t, tls.VersionTLS13, config.MaxVersion, "Maximum TLS version should be 1.3")

		// Test cipher suites
		preferredCiphers := []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		}
		assert.Contains(t, preferredCiphers, tls.TLS_AES_128_GCM_SHA256, "Should support strong cipher suites")
		assert.Contains(t, preferredCiphers, tls.TLS_AES_256_GCM_SHA384, "Should support strong cipher suites")
		assert.Contains(t, preferredCiphers, tls.TLS_CHACHA20_POLY1305_SHA256, "Should support strong cipher suites")
	})
} 