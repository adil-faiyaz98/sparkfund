package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	CommitSHA string `json:"commit_sha,omitempty"`
}

// HealthHandler handles health check requests
type HealthHandler struct {
	version   string
	commitSHA string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(version, commitSHA string) *HealthHandler {
	return &HealthHandler{
		version:   version,
		commitSHA: commitSHA,
	}
}

// RegisterRoutes registers the health check routes
func (h *HealthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/health", h.HealthCheck)
}

// HealthCheck handles health check requests
// @Summary Health check
// @Description Check the health of the service
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "healthy",
		Version:   h.version,
		CommitSHA: h.commitSHA,
	})
}
