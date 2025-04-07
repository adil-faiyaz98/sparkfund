package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// HealthResponse represents the response from the health check endpoint
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	CommitSHA string            `json:"commitSha"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]Check  `json:"checks"`
}

// Check represents the status of a single health check
type Check struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// HealthChecker defines the interface for health checks
type HealthChecker interface {
	Check() (bool, string)
}

// HealthHandler handles health check endpoints
type HealthHandler struct {
	version   string
	commitSHA string
	checks    map[string]HealthChecker
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version, commitSHA string) *HealthHandler {
	return &HealthHandler{
		version:   version,
		commitSHA: commitSHA,
		checks:    make(map[string]HealthChecker),
	}
}

// AddCheck adds a health check
func (h *HealthHandler) AddCheck(name string, checker HealthChecker) {
	h.checks[name] = checker
}

// HealthCheck handles the /health endpoint
func (h *HealthHandler) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		checks := make(map[string]Check)
		status := "healthy"

		for name, checker := range h.checks {
			isHealthy, message := checker.Check()
			checkStatus := "healthy"
			if !isHealthy {
				status = "unhealthy"
				checkStatus = "unhealthy"
			}
			checks[name] = Check{
				Status:  checkStatus,
				Message: message,
			}
		}

		response := HealthResponse{
			Status:    status,
			Version:   h.version,
			CommitSHA: h.commitSHA,
			Timestamp: time.Now(),
			Uptime:    time.Since(startTime).String(),
			Checks:    checks,
		}

		statusCode := http.StatusOK
		if status != "healthy" {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, response)
	}
}

// LivenessCheck handles the /live endpoint
func (h *HealthHandler) LivenessCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	}
}

// ReadinessCheck handles the /ready endpoint
func (h *HealthHandler) ReadinessCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		allReady := true
		for _, checker := range h.checks {
			if ready, _ := checker.Check(); !ready {
				allReady = false
				break
			}
		}

		if allReady {
			c.JSON(http.StatusOK, gin.H{"status": "ready"})
			return
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
	}
}

// DatabaseChecker checks database health
type DatabaseChecker struct {
	Check func() (bool, string)
}

// Check implements the HealthChecker interface
func (d *DatabaseChecker) Check() (bool, string) {
	return d.Check()
}

// CacheChecker checks cache health
type CacheChecker struct {
	Check func() (bool, string)
}

// Check implements the HealthChecker interface
func (c *CacheChecker) Check() (bool, string) {
	return c.Check()
}

// ExternalServiceChecker checks external service health
type ExternalServiceChecker struct {
	Name  string
	URL   string
	Check func() (bool, string)
}

// Check implements the HealthChecker interface
func (e *ExternalServiceChecker) Check() (bool, string) {
	return e.Check()
}
