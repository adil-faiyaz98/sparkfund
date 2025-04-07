package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIVersion represents an API version
type APIVersion string

const (
	// V1 is the first version of the API
	V1 APIVersion = "v1"
	
	// V2 is the second version of the API
	V2 APIVersion = "v2"
	
	// Latest is the latest version of the API
	Latest APIVersion = V1
)

// VersionHeader is the HTTP header used for API versioning
const VersionHeader = "Accept-Version"

// VersionMiddleware handles API versioning
func VersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if version is specified in the URL path
		path := c.Request.URL.Path
		if strings.Contains(path, "/api/v") {
			// Version is already in the path, no need to modify
			c.Next()
			return
		}
		
		// Check if version is specified in the header
		version := c.GetHeader(VersionHeader)
		if version == "" {
			// Default to latest version
			version = string(Latest)
		}
		
		// Validate version
		switch APIVersion(version) {
		case V1, V2:
			// Valid version
		default:
			// Invalid version, default to latest
			version = string(Latest)
		}
		
		// Store version in context
		c.Set("api_version", version)
		
		// Modify path to include version
		if strings.HasPrefix(path, "/api/") {
			newPath := "/api/" + version + path[5:]
			c.Request.URL.Path = newPath
		}
		
		c.Next()
	}
}

// GetVersion returns the API version from the context
func GetVersion(c *gin.Context) APIVersion {
	version, exists := c.Get("api_version")
	if !exists {
		return Latest
	}
	return APIVersion(version.(string))
}

// VersionedHandler handles different API versions
type VersionedHandler struct {
	handlers map[APIVersion]gin.HandlerFunc
}

// NewVersionedHandler creates a new versioned handler
func NewVersionedHandler() *VersionedHandler {
	return &VersionedHandler{
		handlers: make(map[APIVersion]gin.HandlerFunc),
	}
}

// AddHandler adds a handler for a specific version
func (vh *VersionedHandler) AddHandler(version APIVersion, handler gin.HandlerFunc) {
	vh.handlers[version] = handler
}

// Handle handles the request based on the API version
func (vh *VersionedHandler) Handle(c *gin.Context) {
	version := GetVersion(c)
	
	// Check if handler exists for this version
	handler, exists := vh.handlers[version]
	if !exists {
		// Try to use the latest version
		handler, exists = vh.handlers[Latest]
		if !exists {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "API version not implemented",
			})
			c.Abort()
			return
		}
	}
	
	handler(c)
}
