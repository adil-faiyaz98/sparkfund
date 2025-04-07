package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/adil-faiyaz98/sparkfund/services/service-template/internal/api/routes"
	"github.com/adil-faiyaz98/sparkfund/services/service-template/internal/config"
	"github.com/adil-faiyaz98/sparkfund/services/service-template/internal/database"
)

// Server represents the HTTP server
type Server struct {
	config     *config.Config
	router     *gin.Engine
	httpServer *http.Server
	logger     *logrus.Logger
	db         *database.Database
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config) (*Server, error) {
	// Set Gin mode
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Create logger
	logger := logrus.New()
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %v", err)
	}
	logger.SetLevel(level)

	if cfg.Log.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	// Connect to database
	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Create server
	server := &Server{
		config:     cfg,
		router:     router,
		httpServer: httpServer,
		logger:     logger,
		db:         db,
	}

	// Set up middleware and routes
	server.setupMiddleware()
	server.setupRoutes()

	return server, nil
}

// setupMiddleware sets up middleware for the server
func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.router.Use(gin.Recovery())

	// Logger middleware
	s.router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// setupRoutes sets up routes for the server
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Metrics
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := s.router.Group("/api/v1")
	routes.SetupRoutes(api, s.db, s.logger)
}

// Start starts the server
func (s *Server) Start() error {
	s.logger.Infof("Starting server on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	// Close database connection
	if err := s.db.Close(); err != nil {
		s.logger.Errorf("Failed to close database connection: %v", err)
	}

	// Shutdown HTTP server
	return s.httpServer.Shutdown(ctx)
}
