package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/kyc-service/internal/config"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	config     *config.Config
}

// New creates a new server instance
func New(cfg *config.Config, router *gin.Engine) (*Server, error) {
	return &Server{
		router: router,
		config: cfg,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf(":%s", s.config.Server.Port),
		Handler:        s.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Start server
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// RegisterRoutes registers all routes
func (s *Server) RegisterRoutes() {
	// API version group
	v1 := s.router.Group("/api/v1")
	{
		// KYC routes
		kyc := v1.Group("/kyc")
		{
			kyc.POST("/", s.handleCreateKYC)
			kyc.GET("/:userID", s.handleGetKYC)
			kyc.PUT("/", s.handleUpdateKYC)
			kyc.PATCH("/:userID/status", s.handleUpdateKYCStatus)
			kyc.GET("/", s.handleListKYC)
			kyc.DELETE("/:userID", s.handleDeleteKYC)
		}

		// Health check
		v1.GET("/health", s.handleHealth)
	}

	// Metrics endpoint
	s.router.GET("/metrics", s.handleMetrics)
}

// handleCreateKYC handles KYC creation
func (s *Server) handleCreateKYC(c *gin.Context) {
	// TODO: Implement KYC creation
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// handleGetKYC handles getting KYC by user ID
func (s *Server) handleGetKYC(c *gin.Context) {
	// TODO: Implement KYC retrieval
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// handleUpdateKYC handles KYC update
func (s *Server) handleUpdateKYC(c *gin.Context) {
	// TODO: Implement KYC update
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// handleUpdateKYCStatus handles KYC status update
func (s *Server) handleUpdateKYCStatus(c *gin.Context) {
	// TODO: Implement KYC status update
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// handleListKYC handles listing KYC records
func (s *Server) handleListKYC(c *gin.Context) {
	// TODO: Implement KYC listing
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// handleDeleteKYC handles KYC deletion
func (s *Server) handleDeleteKYC(c *gin.Context) {
	// TODO: Implement KYC deletion
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// handleHealth handles health check requests
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// handleMetrics handles metrics requests
func (s *Server) handleMetrics(c *gin.Context) {
	// Prometheus metrics handler is already registered via middleware
	c.Next()
}
