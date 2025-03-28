package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/credit-scoring-service/internal/errors"
	"go.uber.org/zap"
)

type HealthHandler struct {
	logger *zap.Logger
	db     interface {
		Ping() error
	}
	redis interface {
		Ping() error
	}
}

func NewHealthHandler(logger *zap.Logger, db, redis interface{}) *HealthHandler {
	return &HealthHandler{
		logger: logger,
		db:     db,
		redis:  redis,
	}
}

func (h *HealthHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", h.CheckHealth)
	router.GET("/readiness", h.CheckReadiness)
	router.GET("/liveness", h.CheckLiveness)
}

type HealthResponse struct {
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Components  map[string]Status `json:"components"`
	Version     string            `json:"version"`
	Environment string            `json:"environment"`
}

type Status struct {
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
	LatencyMS int64  `json:"latency_ms,omitempty"`
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	start := time.Now()
	response := &HealthResponse{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Components:  make(map[string]Status),
		Version:     "1.0.0",
		Environment: "production",
	}

	// Check database
	dbStart := time.Now()
	if err := h.db.Ping(); err != nil {
		response.Status = "degraded"
		response.Components["database"] = Status{
			Status:    "unhealthy",
			Message:   err.Error(),
			LatencyMS: time.Since(dbStart).Milliseconds(),
		}
	} else {
		response.Components["database"] = Status{
			Status:    "healthy",
			LatencyMS: time.Since(dbStart).Milliseconds(),
		}
	}

	// Check Redis
	redisStart := time.Now()
	if err := h.redis.Ping(); err != nil {
		response.Status = "degraded"
		response.Components["redis"] = Status{
			Status:    "unhealthy",
			Message:   err.Error(),
			LatencyMS: time.Since(redisStart).Milliseconds(),
		}
	} else {
		response.Components["redis"] = Status{
			Status:    "healthy",
			LatencyMS: time.Since(redisStart).Milliseconds(),
		}
	}

	// Set overall status
	if response.Status == "degraded" {
		c.JSON(http.StatusServiceUnavailable, response)
	} else {
		c.JSON(http.StatusOK, response)
	}

	h.logger.Info("health check completed",
		zap.String("status", response.Status),
		zap.Duration("duration", time.Since(start)),
	)
}

func (h *HealthHandler) CheckReadiness(c *gin.Context) {
	response := &HealthResponse{
		Status:      "ready",
		Timestamp:   time.Now(),
		Components:  make(map[string]Status),
		Version:     "1.0.0",
		Environment: "production",
	}

	// Check database
	if err := h.db.Ping(); err != nil {
		response.Status = "not_ready"
		response.Components["database"] = Status{
			Status:  "not_ready",
			Message: err.Error(),
		}
	} else {
		response.Components["database"] = Status{
			Status: "ready",
		}
	}

	// Check Redis
	if err := h.redis.Ping(); err != nil {
		response.Status = "not_ready"
		response.Components["redis"] = Status{
			Status:  "not_ready",
			Message: err.Error(),
		}
	} else {
		response.Components["redis"] = Status{
			Status: "ready",
		}
	}

	if response.Status == "not_ready" {
		c.JSON(http.StatusServiceUnavailable, response)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func (h *HealthHandler) CheckLiveness(c *gin.Context) {
	response := &HealthResponse{
		Status:      "alive",
		Timestamp:   time.Now(),
		Components:  make(map[string]Status),
		Version:     "1.0.0",
		Environment: "production",
	}

	// Basic liveness check
	response.Components["service"] = Status{
		Status: "alive",
	}

	c.JSON(http.StatusOK, response)
} 