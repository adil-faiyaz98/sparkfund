package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sparkfund/kyc-service/internal/middleware"
	"github.com/sparkfund/kyc-service/internal/service"
)

type Router struct {
	engine     *gin.Engine
	handlers   []interface{ RegisterRoutes(*gin.RouterGroup) }
	middleware *middleware.Middleware
}

func NewRouter(services *service.Services, version string) *Router {
	base := NewBaseHandler(services, version)
	mid := middleware.New(services)

	return &Router{
		engine:     gin.New(),
		middleware: mid,
		handlers: []interface{ RegisterRoutes(*gin.RouterGroup) }{
			NewHealthHandler(base),
			NewKYCHandler(base),
			NewMetricsHandler(base),
		},
	}
}

func (r *Router) Setup() {
	// Global middleware
	r.engine.Use(
		gin.Recovery(),
		gin.Logger(),
		r.middleware.RequestID(),
		r.middleware.Cors(),
	)

	// API routes
	v1 := r.engine.Group("/v1")
	v1.Use(r.middleware.Auth())

	for _, h := range r.handlers {
		h.RegisterRoutes(v1)
	}
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}
