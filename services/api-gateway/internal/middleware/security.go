package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// SecurityConfig holds all security-related configuration
type SecurityConfig struct {
	RateLimit struct {
		RequestsPerMinute int
		BurstSize         int
	}
	DoS struct {
		MaxHeaderSize    int64
		MaxBodySize      int64
		MaxConnections   int
		ConnectionWindow time.Duration
	}
	TLS struct {
		MinVersion   uint16
		CipherSuites []uint16
		CertFile     string
		KeyFile      string
	}
	JWT struct {
		SecretKey     []byte
		TokenExpiry   time.Duration
		RefreshExpiry time.Duration
	}
	FileUpload struct {
		MaxSize      int64
		AllowedTypes []string
		MaxFiles     int
	}
}

// SecurityMiddleware implements all security features
type SecurityMiddleware struct {
	config      SecurityConfig
	rateLimiter *RateLimiter
	connTracker *ConnectionTracker
	ipWhitelist map[string]bool
	mu          sync.RWMutex
}

// NewSecurityMiddleware creates a new security middleware instance
func NewSecurityMiddleware(config SecurityConfig) *SecurityMiddleware {
	return &SecurityMiddleware{
		config:      config,
		rateLimiter: NewRateLimiter(config.RateLimit.RequestsPerMinute, config.RateLimit.BurstSize),
		connTracker: NewConnectionTracker(config.DoS.MaxConnections, config.DoS.ConnectionWindow),
		ipWhitelist: make(map[string]bool),
	}
}

// Apply applies all security middleware to the Gin router
func (sm *SecurityMiddleware) Apply(router *gin.Engine) {
	// Global middleware
	router.Use(sm.IPValidation())
	router.Use(sm.RateLimiting())
	router.Use(sm.DoSProtection())
	router.Use(sm.TLSEnforcement())
	router.Use(sm.JWTValidation())
	router.Use(sm.PathTraversalProtection())
	router.Use(sm.FileUploadProtection())
	router.Use(sm.RansomwareProtection())
	router.Use(sm.LogInjectionProtection())
}

// IPValidation middleware validates and sanitizes IP addresses
func (sm *SecurityMiddleware) IPValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get real IP from trusted headers
		realIP := c.GetHeader("X-Real-IP")
		if realIP == "" {
			realIP = c.GetHeader("X-Forwarded-For")
		}

		// Validate IP format and check whitelist
		if !sm.isValidIP(realIP) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid IP address"})
			return
		}

		// Set validated IP in context
		c.Set("client_ip", realIP)
		c.Next()
	}
}

// RateLimiting middleware implements rate limiting
func (sm *SecurityMiddleware) RateLimiting() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.GetString("client_ip")
		if !sm.rateLimiter.Allow(clientIP) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}
		c.Next()
	}
}

// DoSProtection middleware implements DoS protection
func (sm *SecurityMiddleware) DoSProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check header size
		if c.Request.Header.Get("Content-Length") != "" {
			size := c.Request.ContentLength
			if size > sm.config.DoS.MaxHeaderSize {
				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Header too large"})
				return
			}
		}

		// Check body size
		if c.Request.Body != nil {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}
			if int64(len(body)) > sm.config.DoS.MaxBodySize {
				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Body too large"})
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// Track connections
		if !sm.connTracker.Allow(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many connections"})
			return
		}

		c.Next()
	}
}

// TLSEnforcement middleware enforces TLS requirements
func (sm *SecurityMiddleware) TLSEnforcement() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !c.Request.URL.IsAbs() || c.Request.URL.Scheme != "https" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "HTTPS required"})
			return
		}

		// Check TLS version and cipher suite from request
		if tlsConn := c.Request.TLS; tlsConn != nil {
			if tlsConn.Version < sm.config.TLS.MinVersion {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient TLS version"})
				return
			}
			if !contains(sm.config.TLS.CipherSuites, tlsConn.CipherSuite) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Weak cipher suite"})
				return
			}
		}

		c.Next()
	}
}

// JWTValidation middleware validates JWT tokens
func (sm *SecurityMiddleware) JWTValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		// Remove "Bearer " prefix
		token = strings.TrimPrefix(token, "Bearer ")

		// Parse and validate token
		claims := &jwt.StandardClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return sm.config.JWT.SecretKey, nil
		})

		if err != nil || !parsedToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Check token expiry
		if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}

		// Set user claims in context
		c.Set("user_id", claims.Subject)
		c.Next()
	}
}

// PathTraversalProtection middleware prevents path traversal attacks
func (sm *SecurityMiddleware) PathTraversalProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Check for path traversal attempts
		if containsPathTraversal(path) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid path"})
			return
		}

		// Sanitize path
		sanitizedPath := filepath.Clean(path)
		if sanitizedPath != path {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid path"})
			return
		}

		c.Next()
	}
}

// FileUploadProtection middleware protects against malicious file uploads
func (sm *SecurityMiddleware) FileUploadProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

		// Check content type
		contentType := c.GetHeader("Content-Type")
		if !strings.Contains(contentType, "multipart/form-data") {
			c.Next()
			return
		}

		// Parse multipart form
		if err := c.Request.ParseMultipartForm(sm.config.FileUpload.MaxSize); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		// Check number of files
		if len(c.Request.MultipartForm.File) > sm.config.FileUpload.MaxFiles {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Too many files"})
			return
		}

		// Validate each file
		for _, files := range c.Request.MultipartForm.File {
			for _, file := range files {
				if !isValidFileType(file.Filename, sm.config.FileUpload.AllowedTypes) {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
					return
				}
			}
		}

		c.Next()
	}
}

// RansomwareProtection middleware detects and prevents ransomware attacks
func (sm *SecurityMiddleware) RansomwareProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.Next()
			return
		}

		// Check for suspicious file operations
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			if len(bodyBytes) > 0 {
				var body map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &body); err == nil {
					if isSuspiciousOperation(body) {
						c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Suspicious operation detected"})
						return
					}
				}
			}

			// Restore the body for downstream handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		c.Next()
	}
}

// LogInjectionProtection middleware prevents log injection attacks
func (sm *SecurityMiddleware) LogInjectionProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize headers
		for key, values := range c.Request.Header {
			for i, value := range values {
				c.Request.Header[key][i] = sanitizeLogValue(value)
			}
		}

		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				c.Request.URL.Query()[key][i] = sanitizeLogValue(value)
			}
		}

		c.Next()
	}
}

// AddToIPWhitelist adds an IP address to the whitelist
func (sm *SecurityMiddleware) AddToIPWhitelist(ip string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.ipWhitelist[ip] = true
}

// Helper functions
func (sm *SecurityMiddleware) isValidIP(ip string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.ipWhitelist[ip]
}

func contains(slice []uint16, item uint16) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsPathTraversal(path string) bool {
	traversalPatterns := []string{
		"../",
		"..\\",
		"%2e%2e%2f",
		"%252e%252e%252f",
		"..%252f",
		"%252e%252e%252f",
		"....//",
		"..%2f",
		"%2e%2e%2f",
	}

	for _, pattern := range traversalPatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func isValidFileType(filename string, allowedTypes []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

func isSuspiciousOperation(body map[string]interface{}) bool {
	// Check for bulk operations
	if files, ok := body["files"].([]interface{}); ok && len(files) > 100 {
		return true
	}

	// Check for encryption operations
	if operation, ok := body["operation"].(string); ok {
		if operation == "encrypt" || operation == "delete" {
			return true
		}
	}

	// Check for suspicious file patterns
	if files, ok := body["files"].([]interface{}); ok {
		for _, file := range files {
			if str, ok := file.(string); ok {
				if strings.Contains(str, "*") || strings.Contains(str, "..") {
					return true
				}
			}
		}
	}

	return false
}

func sanitizeLogValue(value string) string {
	// Remove control characters
	value = strings.Map(func(r rune) rune {
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			return -1
		}
		return r
	}, value)

	// Escape special characters
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "\t", "\\t")

	return value
}

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	requests map[string]*tokenBucket
	mu       sync.RWMutex
	rate     int
	burst    int
}

type tokenBucket struct {
	tokens     int
	lastUpdate time.Time
}

func NewRateLimiter(rate, burst int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]*tokenBucket),
		rate:     rate,
		burst:    burst,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.requests[ip]
	if !exists {
		bucket = &tokenBucket{
			tokens:     rl.burst,
			lastUpdate: time.Now(),
		}
		rl.requests[ip] = bucket
	}

	now := time.Now()
	elapsed := now.Sub(bucket.lastUpdate)
	tokensToAdd := int(elapsed.Seconds() * float64(rl.rate))
	if tokensToAdd > 0 {
		bucket.tokens = min(bucket.tokens+tokensToAdd, rl.burst)
		bucket.lastUpdate = now
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

// ConnectionTracker implements connection limiting
type ConnectionTracker struct {
	connections map[string][]time.Time
	mu          sync.RWMutex
	maxConns    int
	window      time.Duration
}

func NewConnectionTracker(maxConns int, window time.Duration) *ConnectionTracker {
	return &ConnectionTracker{
		connections: make(map[string][]time.Time),
		maxConns:    maxConns,
		window:      window,
	}
}

func (ct *ConnectionTracker) Allow(ip string) bool {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	now := time.Now()
	conns := ct.connections[ip]

	// Remove old connections
	valid := conns[:0]
	for _, t := range conns {
		if now.Sub(t) <= ct.window {
			valid = append(valid, t)
		}
	}

	if len(valid) >= ct.maxConns {
		return false
	}

	valid = append(valid, now)
	ct.connections[ip] = valid
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
