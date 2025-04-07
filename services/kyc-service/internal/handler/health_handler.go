package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HealthHandler struct {
	*BaseHandler
	startTime time.Time
}

func NewHealthHandler(base *BaseHandler) *HealthHandler {
	return &HealthHandler{
		BaseHandler: base,
		startTime:   time.Now(),
	}
}

func (h *HealthHandler) RegisterRoutes(r *gin.RouterGroup) {
	health := r.Group("/health")
	{
		health.GET("", h.HealthCheck)
		health.GET("/ready", h.ReadinessCheck)
		health.GET("/live", h.LivenessCheck)
		health.GET("/status", h.DetailedStatus)
	}
}

func (h *HealthHandler) DetailedStatus(c *gin.Context) {
	status := h.services.Health.GetDetailedStatus(c.Request.Context())
	c.JSON(http.StatusOK, status)
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	checks := h.healthService.PerformHealthChecks()

	status := "healthy"
	statusCode := http.StatusOK

	for _, checkStatus := range checks {
		if checkStatus != "healthy" {
			status = "unhealthy"
			statusCode = http.StatusServiceUnavailable
			break
		}
	}

	response := HealthResponse{
		Status:    status,
		Version:   h.version,
		Uptime:    time.Since(h.startTime).String(),
		Timestamp: time.Now(),
		Checks:    checks,
	}

	c.JSON(statusCode, response)
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	if !h.healthService.IsReady() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}
