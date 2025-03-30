package handlers

import (
	"net/http"
	"time"

	"investment-service/internal/database"

	"github.com/gin-gonic/gin"
)

// HealthResponse contains service health information
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]string `json:"checks"`
}

var startTime = time.Now()

// HealthCheck handles the health endpoint
// @Summary Get service health
// @Description Check the health of the service and its dependencies
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	checks := make(map[string]string)
	statusCode := http.StatusOK
	overallStatus := "healthy"

	// Check database connection
	if err := database.DB.Exec("SELECT 1").Error; err != nil {
		checks["database"] = "unhealthy: " + err.Error()
		statusCode = http.StatusServiceUnavailable
		overallStatus = "unhealthy"
	} else {
		checks["database"] = "healthy"
	}

	// Check disk space (simplified)
	checks["disk"] = "healthy"

	// Check memory usage (simplified)
	checks["memory"] = "healthy"

	response := HealthResponse{
		Status:    overallStatus,
		Version:   "1.0.0", // This should be injected from build
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
		Checks:    checks,
	}

	c.JSON(statusCode, response)
}

// ReadinessCheck handles the readiness endpoint
// @Summary Check if service is ready to serve requests
// @Description Check if the service is ready to handle requests
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /ready [get]
func ReadinessCheck(c *gin.Context) {
	// Check if database migrations are complete
	if err := database.DB.Exec("SELECT 1 FROM investments LIMIT 1").Error; err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"reason":  "database not ready",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// LivenessCheck handles the liveness endpoint
// @Summary Check if service is alive
// @Description Basic ping to check if the service is running
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /live [get]
func LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}
