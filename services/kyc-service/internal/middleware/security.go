package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityConfig holds security middleware configuration
type SecurityConfig struct {
	EnableHSTS          bool
	EnableXSSProtection bool
	EnableCSRF          bool
	EnableCORS          bool
	AllowedOrigins      []string
	AllowedMethods      []string
	AllowedHeaders      []string
	MaxAge              int
	ContentTypeOptions  bool
	FrameOptions        string
	ReferrerPolicy      string
	PermissionsPolicy   map[string][]string
}

// DefaultSecurityConfig returns default security configuration
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		EnableHSTS:          true,
		EnableXSSProtection: true,
		EnableCSRF:          true,
		EnableCORS:          true,
		AllowedOrigins:      []string{"*"},
		AllowedMethods:      []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:      []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:              86400,
		ContentTypeOptions:  true,
		FrameOptions:        "DENY",
		ReferrerPolicy:      "strict-origin-when-cross-origin",
		PermissionsPolicy: map[string][]string{
			"geolocation": {},
			"camera":      {},
			"microphone":  {},
		},
	}
}

// SecurityMiddleware provides security headers and protection against various attack vectors
func SecurityMiddleware(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// HSTS
		if config.EnableHSTS {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// XSS Protection
		if config.EnableXSSProtection {
			c.Header("X-XSS-Protection", "1; mode=block")
		}

		// Content Type Options
		if config.ContentTypeOptions {
			c.Header("X-Content-Type-Options", "nosniff")
		}

		// Frame Options
		if config.FrameOptions != "" {
			c.Header("X-Frame-Options", config.FrameOptions)
		}

		// Referrer Policy
		if config.ReferrerPolicy != "" {
			c.Header("Referrer-Policy", config.ReferrerPolicy)
		}

		// Permissions Policy
		if len(config.PermissionsPolicy) > 0 {
			var policy []string
			for feature, allowlist := range config.PermissionsPolicy {
				if len(allowlist) == 0 {
					policy = append(policy, feature+"=()")
				} else {
					policy = append(policy, feature+"=("+strings.Join(allowlist, " ")+")")
				}
			}
			c.Header("Permissions-Policy", strings.Join(policy, ", "))
		}

		// CORS
		if config.EnableCORS {
			origin := c.GetHeader("Origin")
			if origin != "" {
				allowed := false
				for _, allowedOrigin := range config.AllowedOrigins {
					if allowedOrigin == "*" || allowedOrigin == origin {
						allowed = true
						break
					}
				}
				if allowed {
					c.Header("Access-Control-Allow-Origin", origin)
					c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
					c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
					c.Header("Access-Control-Max-Age", string(config.MaxAge))
					if c.Request.Method == "OPTIONS" {
						c.AbortWithStatus(http.StatusNoContent)
						return
					}
				}
			}
		}

		// CSRF Protection
		if config.EnableCSRF {
			if c.Request.Method != "GET" && c.Request.Method != "HEAD" && c.Request.Method != "OPTIONS" {
				token := c.GetHeader("X-CSRF-Token")
				if token == "" {
					c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
					c.Abort()
					return
				}
				// Validate CSRF token here
				// This is a simplified example. In production, you should use a proper CSRF token validation
			}
		}

		// Additional Security Headers
		c.Header("X-Requested-With", "XMLHttpRequest")
		c.Header("X-Download-Options", "noopen")
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		c.Next()
	}
}

// ValidateRequestSize validates the request size
func ValidateRequestSize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Request too large"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ValidateContentType validates the content type
func ValidateContentType(allowedTypes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		allowed := false
		for _, t := range allowedTypes {
			if strings.Contains(contentType, t) {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported media type"})
			c.Abort()
			return
		}
		c.Next()
	}
}
