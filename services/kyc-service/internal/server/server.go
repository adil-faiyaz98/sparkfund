package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sparkfund/shared/config"
	"github.com/sparkfund/shared/handlers/health"
)

// Version is the service version
var Version = "1.0.0"

// CommitSHA is the git commit SHA
var CommitSHA = "unknown"

// DatabaseHealthChecker checks database health
type DatabaseHealthChecker struct{}

// Check implements the health.HealthChecker interface
func (d *DatabaseHealthChecker) Check() (bool, string) {
	// TODO: Implement actual database health check
	return true, "Database is healthy"
}

// CacheHealthChecker checks cache health
type CacheHealthChecker struct{}

// Check implements the health.HealthChecker interface
func (c *CacheHealthChecker) Check() (bool, string) {
	// TODO: Implement actual cache health check
	return true, "Cache is healthy"
}

type Server struct {
	router *gin.Engine
	config *config.BaseConfig
	health *health.HealthHandler
}

func NewServer(cfg *config.BaseConfig) *Server {
	router := gin.Default()

	healthHandler := health.NewHealthHandler(
		Version,
		CommitSHA,
	)

	// Add health checks
	healthHandler.AddCheck("database", &DatabaseHealthChecker{})
	healthHandler.AddCheck("cache", &CacheHealthChecker{})

	return &Server{
		router: router,
		config: cfg,
		health: healthHandler,
	}
}

func (s *Server) SetupRoutes() {
	// Health endpoints
	s.router.GET("/health", s.health.HealthCheck())
	s.router.GET("/live", s.health.LivenessCheck())
	s.router.GET("/ready", s.health.ReadinessCheck())

	// Service-specific routes
	// ...
}
