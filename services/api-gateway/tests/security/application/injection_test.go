package application

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/api-gateway/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestInjectionAttacks(t *testing.T) {
	// Create test router
	router := gin.New()
	sm := middleware.NewSecurityMiddleware(middleware.SecurityConfig{})
	sm.Apply(router)

	// Test SQL Injection
	t.Run("SQL Injection", func(t *testing.T) {
		// Test basic SQL injection
		sqlPayloads := []string{
			"' OR '1'='1",
			"' UNION SELECT * FROM users; --",
			"'; DROP TABLE users; --",
			"' OR '1'='1' --",
			"' OR 'x'='x",
			"' OR '1'='1' #",
			"' OR '1'='1' /*",
			"admin' --",
			"admin' #",
			"admin'/*",
		}

		for _, payload := range sqlPayloads {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/login", bytes.NewBufferString(`{"username": "`+payload+`"}`))
			router.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "SQL injection attempt should be blocked: %s", payload)
		}

		// Test blind SQL injection
		blindPayloads := []string{
			"' AND 1=1 --",
			"' AND 1=2 --",
			"' AND SLEEP(5) --",
			"' AND (SELECT * FROM (SELECT(SLEEP(5)))a) --",
		}

		for _, payload := range blindPayloads {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/search", bytes.NewBufferString(`{"query": "`+payload+`"}`))
			router.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "Blind SQL injection attempt should be blocked: %s", payload)
		}

		// Test time-based SQL injection
		timePayloads := []string{
			"' AND (SELECT * FROM (SELECT(SLEEP(5)))a) --",
			"' AND (SELECT * FROM (SELECT(SLEEP(5)))a) #",
			"' AND (SELECT * FROM (SELECT(SLEEP(5)))a) /*",
		}

		for _, payload := range timePayloads {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/data", bytes.NewBufferString(`{"id": "`+payload+`"}`))
			router.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "Time-based SQL injection attempt should be blocked: %s", payload)
		}
	})

	// Test XSS Attacks
	t.Run("XSS Attacks", func(t *testing.T) {
		// Test basic XSS
		xssPayloads := []string{
			"<script>alert('xss')</script>",
			"<img src=x onerror=alert('xss')>",
			"<svg><script>alert('xss')</script></svg>",
			"<body onload=alert('xss')>",
			"<input autofocus onfocus=alert('xss')>",
			"<select autofocus onfocus=alert('xss')>",
			"<textarea autofocus onfocus=alert('xss')>",
			"<keygen autofocus onfocus=alert('xss')>",
			"<div/onmouseover='alert(1)'>style=width:100%;height:100%;position:fixed;left:0;top:0",
			"<svg><script>alert('xss')</script></svg>",
		}

		for _, payload := range xssPayloads {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/comment", bytes.NewBufferString(`{"text": "`+payload+`"}`))
			router.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "XSS attempt should be blocked: %s", payload)
		}

		// Test encoded XSS
		encodedPayloads := []string{
			"&#x3C;script&#x3E;alert('xss')&#x3C;/script&#x3E;",
			"&#60;script&#62;alert('xss')&#60;/script&#62;",
			"&lt;script&gt;alert('xss')&lt;/script&gt;",
			"%3Cscript%3Ealert('xss')%3C/script%3E",
			"\\x3Cscript\\x3Ealert('xss')\\x3C/script\\x3E",
		}

		for _, payload := range encodedPayloads {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/comment", bytes.NewBufferString(`{"text": "`+payload+`"}`))
			router.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "Encoded XSS attempt should be blocked: %s", payload)
		}

		// Test DOM-based XSS
		domPayloads := []string{
			"javascript:alert('xss')",
			"data:text/html,<script>alert('xss')</script>",
			"vbscript:alert('xss')",
			"data:application/x-javascript,alert('xss')",
		}

		for _, payload := range domPayloads {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/url", bytes.NewBufferString(`{"url": "`+payload+`"}`))
			router.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "DOM-based XSS attempt should be blocked: %s", payload)
		}

		// Test event handler XSS
		eventPayloads := []string{
			"onerror=alert('xss')",
			"onload=alert('xss')",
			"onmouseover=alert('xss')",
			"onmouseout=alert('xss')",
			"onclick=alert('xss')",
			"onkeypress=alert('xss')",
			"onkeydown=alert('xss')",
			"onkeyup=alert('xss')",
			"onfocus=alert('xss')",
			"onblur=alert('xss')",
		}

		for _, payload := range eventPayloads {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/event", bytes.NewBufferString(`{"handler": "`+payload+`"}`))
			router.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusOK, w.Code, "Event handler XSS attempt should be blocked: %s", payload)
		}
	})
} 