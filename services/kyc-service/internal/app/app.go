package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sparkfund/services/kyc-service/internal/api"
	"sparkfund/services/kyc-service/internal/config"
	"sparkfund/services/kyc-service/internal/repository"
	"sparkfund/services/kyc-service/internal/service"
)

// App represents the application
type App struct {
	config     *config.Config
	httpServer *http.Server
	db         *gorm.DB
	services   *service.Services
	router     *api.Router
}

// New creates a new application
func New(cfg *config.Config) (*App, error) {
	// Set gin mode
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Connect to database
	db, err := repository.NewDB(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create repositories
	repos := repository.NewRepositories(db)

	// Create event publisher
	eventPublisher := service.NewEventPublisher(cfg.Events)

	// Create services
	services := service.NewServices(service.ServicesDeps{
		Repos:          repos,
		EventPublisher: eventPublisher,
		Config:         cfg,
	})

	// Create router
	router := api.NewRouter(services, api.RouterConfig{
		Version:   cfg.App.Version,
		CommitSHA: os.Getenv("GIT_COMMIT"),
		Debug:     cfg.App.Environment == "development",
	})

	// Create HTTP server
	httpServer := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:        router.Engine(),
		ReadTimeout:    cfg.Server.Timeout,
		WriteTimeout:   cfg.Server.Timeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return &App{
		config:     cfg,
		httpServer: httpServer,
		db:         db,
		services:   services,
		router:     router,
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	// Start HTTP server
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on %s", a.httpServer.Addr)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
	return nil
}
