package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"sparkfund/services/kyc-service/internal/api/handlers"
	"sparkfund/services/kyc-service/internal/api/middleware"
	"sparkfund/services/kyc-service/internal/service"
)

// Router handles HTTP routing
type Router struct {
	engine *gin.Engine
	config RouterConfig
}

// RouterConfig contains router configuration
type RouterConfig struct {
	Version   string
	CommitSHA string
	Debug     bool
}

// NewRouter creates a new router
func NewRouter(services *service.Services, config RouterConfig) *Router {
	// Create router
	r := &Router{
		engine: gin.New(),
		config: config,
	}

	// Set gin mode
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Add versioning middleware
	r.engine.Use(middleware.VersionMiddleware())

	// Setup middleware
	r.engine.Use(gin.Recovery())
	r.engine.Use(middleware.Logger())
	r.engine.Use(middleware.CORS())

	// Create handlers
	healthHandler := handlers.NewHealthHandler(config.Version, config.CommitSHA)
	documentHandler := handlers.NewDocumentHandler(services.Document)
	kycHandler := handlers.NewKYCHandler(services.KYC)
	verificationHandler := handlers.NewVerificationHandler(services.Verification)

	// Register routes
	api := r.engine.Group("/api/v1")
	{
		// Health check
		healthHandler.RegisterRoutes(api)

		// Metrics
		api.GET("/metrics", gin.WrapH(promhttp.Handler()))

		// Document routes
		documentHandler.RegisterRoutes(api)

		// KYC routes
		kycHandler.RegisterRoutes(api)

		// Verification routes
		verificationHandler.RegisterRoutes(api)
	}

	return r
}

// Engine returns the gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
