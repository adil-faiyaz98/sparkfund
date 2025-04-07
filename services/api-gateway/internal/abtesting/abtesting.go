package abtesting

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Config holds A/B testing configuration
type Config struct {
	Enabled     bool              `mapstructure:"enabled"`
	CookieName  string            `mapstructure:"cookie_name"`
	CookieTTL   int               `mapstructure:"cookie_ttl"`
	Experiments map[string]Test   `mapstructure:"experiments"`
}

// Test represents an A/B test
type Test struct {
	Enabled     bool              `mapstructure:"enabled"`
	Description string            `mapstructure:"description"`
	Variants    map[string]Variant `mapstructure:"variants"`
	DefaultVariant string         `mapstructure:"default_variant"`
	StickySession bool            `mapstructure:"sticky_session"`
	HeaderName  string            `mapstructure:"header_name"`
	QueryParam  string            `mapstructure:"query_param"`
}

// Variant represents a test variant
type Variant struct {
	Weight      int               `mapstructure:"weight"`
	Description string            `mapstructure:"description"`
}

// Manager manages A/B testing
type Manager struct {
	config Config
	logger *logrus.Logger
	mu     sync.RWMutex
	rand   *rand.Rand
}

// NewManager creates a new A/B testing manager
func NewManager(config Config, logger *logrus.Logger) *Manager {
	// Initialize random number generator with seed
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	return &Manager{
		config: config,
		logger: logger,
		rand:   rng,
	}
}

// IsEnabled returns whether A/B testing is enabled
func (m *Manager) IsEnabled() bool {
	return m.config.Enabled
}

// Middleware returns a middleware that handles A/B testing
func (m *Manager) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.config.Enabled {
			c.Next()
			return
		}

		// Process each experiment
		for expName, test := range m.config.Experiments {
			if !test.Enabled {
				continue
			}

			// Determine variant
			variant := m.determineVariant(c, expName, test)

			// Set variant in context
			c.Set("ab_"+expName, variant)

			// Set variant in response header
			c.Header("X-AB-"+expName, variant)

			// If sticky session is enabled, set cookie
			if test.StickySession {
				// Check if cookie already exists
				_, err := c.Cookie("ab_" + expName)
				if err != nil {
					// Set cookie with variant
					c.SetCookie("ab_"+expName, variant, m.config.CookieTTL, "/", "", false, false)
				}
			}

			m.logger.Debugf("A/B test %s: assigned variant %s", expName, variant)
		}

		c.Next()
	}
}

// determineVariant determines the variant for a request
func (m *Manager) determineVariant(c *gin.Context, expName string, test Test) string {
	// Check if variant is specified in header
	if test.HeaderName != "" {
		if variant := c.GetHeader(test.HeaderName); variant != "" {
			if _, ok := test.Variants[variant]; ok {
				return variant
			}
		}
	}

	// Check if variant is specified in query parameter
	if test.QueryParam != "" {
		if variant := c.Query(test.QueryParam); variant != "" {
			if _, ok := test.Variants[variant]; ok {
				return variant
			}
		}
	}

	// Check if sticky session is enabled and cookie exists
	if test.StickySession {
		if variant, err := c.Cookie("ab_" + expName); err == nil {
			if _, ok := test.Variants[variant]; ok {
				return variant
			}
		}
	}

	// Determine variant based on user identifier or random assignment
	userID := getUserIdentifier(c)
	return m.selectVariant(expName, test, userID)
}

// selectVariant selects a variant based on weights and user identifier
func (m *Manager) selectVariant(expName string, test Test, userID string) string {
	// If user ID is available, use it for consistent assignment
	if userID != "" {
		// Create a hash of experiment name and user ID
		hash := sha256.Sum256([]byte(expName + userID))
		hashStr := hex.EncodeToString(hash[:])
		
		// Convert first 8 characters of hash to integer
		hashInt, err := strconv.ParseUint(hashStr[:8], 16, 64)
		if err == nil {
			// Use hash to determine variant
			return m.selectVariantByWeight(test, int(hashInt % 100))
		}
	}

	// Fall back to random assignment
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.selectVariantByWeight(test, m.rand.Intn(100))
}

// selectVariantByWeight selects a variant based on weights
func (m *Manager) selectVariantByWeight(test Test, value int) string {
	// Calculate total weight
	totalWeight := 0
	for _, variant := range test.Variants {
		totalWeight += variant.Weight
	}

	// If total weight is 0, return default variant
	if totalWeight == 0 {
		return test.DefaultVariant
	}

	// Normalize value to total weight
	value = value % totalWeight

	// Select variant based on weight
	currentWeight := 0
	for variantName, variant := range test.Variants {
		currentWeight += variant.Weight
		if value < currentWeight {
			return variantName
		}
	}

	// Fallback to default variant
	return test.DefaultVariant
}

// getUserIdentifier gets a user identifier from the request
func getUserIdentifier(c *gin.Context) string {
	// Try to get user ID from context (set by auth middleware)
	if user, exists := c.Get("user"); exists {
		if userMap, ok := user.(map[string]interface{}); ok {
			if userID, ok := userMap["sub"].(string); ok {
				return userID
			}
		}
	}

	// Try to get from session cookie
	sessionCookie, _ := c.Cookie("session")
	if sessionCookie != "" {
		return sessionCookie
	}

	// Use client IP as fallback
	clientIP := c.ClientIP()
	if clientIP != "" {
		return clientIP
	}

	// Use user agent as last resort
	userAgent := c.Request.UserAgent()
	if userAgent != "" {
		return userAgent
	}

	return ""
}

// GetExperiments returns all experiments
func (m *Manager) GetExperiments() map[string]Test {
	return m.config.Experiments
}

// GetVariant gets the variant for a specific experiment
func (m *Manager) GetVariant(c *gin.Context, expName string) string {
	if variant, exists := c.Get("ab_" + expName); exists {
		if variantStr, ok := variant.(string); ok {
			return variantStr
		}
	}
	
	// Return default variant if not found
	if test, ok := m.config.Experiments[expName]; ok {
		return test.DefaultVariant
	}
	
	return ""
}
