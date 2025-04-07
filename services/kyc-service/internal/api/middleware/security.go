package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/sparkfund/kyc-service/internal/security"
)

type SecurityMiddleware struct {
    config security.SecurityConfig
    threatDetector *security.ThreatDetector
}

func NewSecurityMiddleware(config security.SecurityConfig, threatDetector *security.ThreatDetector) *SecurityMiddleware {
    return &SecurityMiddleware{
        config: config,
        threatDetector: threatDetector,
    }
}

// DetectThreats implements threat detection
func (m *SecurityMiddleware) DetectThreats() gin.HandlerFunc {
    return func(c *gin.Context) {
        score := m.threatDetector.Analyze(c.Request)
        if score > m.config.ThresholdScore {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "error": "Security threat detected",
            })
            return
        }
        c.Next()
    }
}

// RateLimit implements rate limiting
func (m *SecurityMiddleware) RateLimit() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientID := security.GetClientID(c)
        if !security.RateLimiter.Allow(clientID) {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
            })
            return
        }
        c.Next()
    }
}

// SecureHeaders adds security headers
func (m *SecurityMiddleware) SecureHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Next()
    }
}