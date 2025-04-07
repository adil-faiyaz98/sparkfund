package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sparkfund/kyc-service/internal/middleware"
	"github.com/sparkfund/kyc-service/internal/service"
)

type BaseHandler struct {
	services   *service.Services
	version    string
	middleware *middleware.Middleware
}

func NewBaseHandler(services *service.Services, version string) *BaseHandler {
	return &BaseHandler{
		services:   services,
		version:    version,
		middleware: middleware.New(services),
	}
}

// Common handler utilities
func (h *BaseHandler) handleError(c *gin.Context, err error) {
	// Centralized error handling
}

func (h *BaseHandler) validateRequest(c *gin.Context, req interface{}) bool {
	// Common request validation
}
