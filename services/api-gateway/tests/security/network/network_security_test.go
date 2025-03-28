package network

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/api-gateway/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestNetworkSecurity(t *testing.T) {
	// Create test configuration
	config := middleware.NetworkSecurityConfig{
		TCP: struct {
			MaxConnections    int
			ConnectionTimeout time.Duration
			BacklogSize       int
			MaxSynQueue       int
		}{
			MaxConnections:    100,
			ConnectionTimeout: time.Second * 30,
			BacklogSize:       1000,
			MaxSynQueue:       50,
		},
		SSL: struct {
			MinVersion   uint16
			CipherSuites []uint16
			CertFile     string
			KeyFile      string
			StrictSNI    bool
			HSTSEnabled  bool
			HSTSDuration time.Duration
		}{
			MinVersion:   tls.VersionTLS12,
			CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
			HSTSEnabled:  true,
			HSTSDuration: time.Hour * 24 * 365,
		},
		Network: struct {
			AllowedIPRanges []string
			BlockedIPs      []string
			RateLimit       struct {
				RequestsPerSecond int
				BurstSize         int
			}
		}{
			AllowedIPRanges: []string{"192.168.1.0/24"},
			BlockedIPs:      []string{"10.0.0.1"},
			RateLimit: struct {
				RequestsPerSecond int
				BurstSize         int
			}{
				RequestsPerSecond: 100,
				BurstSize:         20,
			},
		},
	}

	// Create test router
	router := gin.New()
	nsm := middleware.NewNetworkSecurityMiddleware(config)
	nsm.Apply(router)

	// Test SSL Stripping Protection
	t.Run("SSL Stripping Protection", func(t *testing.T) {
		// Test HTTP request (should redirect to HTTPS)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusMovedPermanently, w.Code)

		// Test HTTPS request with HSTS
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "https://example.com/test", nil)
		req.TLS = &tls.ConnectionState{
			Version:     tls.VersionTLS12,
			CipherSuite: tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		}
		router.ServeHTTP(w, req)
		assert.Equal(t, "max-age=31536000", w.Header().Get("Strict-Transport-Security"))
	})

	// Test SYN Flood Protection
	t.Run("SYN Flood Protection", func(t *testing.T) {
		clientIP := "192.168.1.1"
		for i := 0; i < 60; i++ {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "https://example.com/test", nil)
			req.RemoteAddr = clientIP
			router.ServeHTTP(w, req)
			if i >= 50 {
				assert.Equal(t, http.StatusTooManyRequests, w.Code, "SYN flood should be detected")
			}
			time.Sleep(time.Millisecond * 10)
		}
	})

	// Test Route Table Protection
	t.Run("Route Table Protection", func(t *testing.T) {
		// Test allowed IP
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "https://example.com/test", nil)
		req.RemoteAddr = "192.168.1.100"
		router.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)

		// Test blocked IP
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "https://example.com/test", nil)
		req.RemoteAddr = "10.0.0.1"
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)

		// Test suspicious routing
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "https://example.com/test", nil)
		req.Header.Set("X-Forwarded-Host", "malicious.com")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	// Test Smurf Attack Protection
	t.Run("Smurf Attack Protection", func(t *testing.T) {
		// Test broadcast IP
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "https://example.com/test", nil)
		req.RemoteAddr = "255.255.255.255"
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)

		// Test ICMP flood
		clientIP := "192.168.1.1"
		for i := 0; i < 100; i++ {
			w = httptest.NewRecorder()
			req = httptest.NewRequest("PING", "https://example.com/test", nil)
			req.RemoteAddr = clientIP
			router.ServeHTTP(w, req)
			if i >= 20 {
				assert.Equal(t, http.StatusTooManyRequests, w.Code, "ICMP flood should be detected")
			}
			time.Sleep(time.Millisecond * 10)
		}
	})

	// Test MAC Address Protection
	t.Run("MAC Address Protection", func(t *testing.T) {
		// Test valid MAC
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "https://example.com/test", nil)
		req.Header.Set("X-Real-MAC", "00:11:22:33:44:55")
		router.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)

		// Test invalid MAC
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "https://example.com/test", nil)
		req.Header.Set("X-Real-MAC", "invalid:mac:address")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	// Test Wireless Security
	t.Run("Wireless Security", func(t *testing.T) {
		// Test WiFi headers
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "https://example.com/test", nil)
		req.Header.Set("X-Wifi-Device", "test-device")
		router.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)

		// Test Bluetooth headers
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "https://example.com/test", nil)
		req.Header.Set("X-Bluetooth-Device", "test-device")
		router.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)

		// Test RFID headers
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "https://example.com/test", nil)
		req.Header.Set("X-RFID-Tag", "test-tag")
		router.ServeHTTP(w, req)
		assert.NotEqual(t, http.StatusForbidden, w.Code)
	})
}
