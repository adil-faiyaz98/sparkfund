package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/sparkfund/credit-scoring-service/internal/config"
	"github.com/sparkfund/credit-scoring-service/internal/database"
	"github.com/sparkfund/credit-scoring-service/internal/handler"
	"github.com/sparkfund/credit-scoring-service/internal/middleware"
	"github.com/sparkfund/credit-scoring-service/internal/repository"
	"github.com/sparkfund/credit-scoring-service/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize logger
	logger := initLogger(cfg)
	defer logger.Sync()

	// Initialize database
	db, err := initDB(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Initialize Redis for rate limiting
	redisClient, err := initRedis(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Initialize dependencies
	creditRepo := repository.NewCreditRepository(db)
	creditService := service.NewCreditService(creditRepo)
	creditHandler := handler.NewCreditHandler(creditService, logger)

	// Initialize health handler
	healthHandler := handler.NewHealthHandler(logger, db, redisClient)

	// Initialize middleware
	securityMiddleware := middleware.NewSecurityMiddleware(cfg, logger)
	monitoringMiddleware := middleware.NewMonitoringMiddleware(logger)

	// Initialize rate limiter
	rateLimiter, err := middleware.NewRateLimiter(&middleware.RateLimiterConfig{
		RequestsPerMinute: cfg.Security.RateLimit.RequestsPerMinute,
		BurstSize:        cfg.Security.RateLimit.BurstSize,
		RedisURL:         cfg.Redis.URL,
		RedisPassword:    cfg.Redis.Password,
		RedisDB:          cfg.Redis.DB,
	}, logger)
	if err != nil {
		logger.Fatal("Failed to initialize rate limiter", zap.Error(err))
	}
	defer rateLimiter.Close()

	// Initialize router
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(monitoringMiddleware.RequestLogger())
	router.Use(monitoringMiddleware.Metrics())
	router.Use(securityMiddleware.RequestID())
	router.Use(securityMiddleware.RequestSizeLimit())
	router.Use(securityMiddleware.ValidateContentType())
	router.Use(securityMiddleware.SecurityHeaders())
	router.Use(securityMiddleware.CORS())
	router.Use(rateLimiter.RateLimit())
	router.Use(middleware.AuthMiddleware(&middleware.JWTConfig{
		SecretKey:     cfg.JWT.SecretKey,
		TokenExpiry:   cfg.JWT.TokenExpiry,
		Issuer:        cfg.JWT.Issuer,
		Audience:      cfg.JWT.Audience,
		AllowedScopes: cfg.JWT.AllowedScopes,
	}, logger))

	// Register routes
	creditHandler.RegisterRoutes(router)
	healthHandler.RegisterRoutes(router)

	// Register metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Start server
	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown server
	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}

func initLogger(cfg *config.Config) *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(getLogLevel(cfg.Logging.Level))
	config.OutputPaths = []string{cfg.Logging.Output}
	config.Encoding = cfg.Logging.Format

	// Configure log rotation
	if cfg.Logging.Output != "stdout" {
		config.OutputPaths = []string{cfg.Logging.Output}
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func initDB(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	config := &gorm.Config{
		Logger: logger.New(
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	return db, nil
}

func initRedis(cfg *config.Config, logger *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            cfg.Redis.URL,
		Password:        cfg.Redis.Password,
		DB:              cfg.Redis.DB,
		PoolSize:        cfg.Redis.PoolSize,
		MinIdleConns:    cfg.Redis.MinIdleConns,
		MaxConnAge:      cfg.Redis.MaxConnAge,
		RequestTimeout:  cfg.Redis.RequestTimeout,
		MaxRetries:      3,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			logger.Info("Connected to Redis")
			return nil
		},
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
} 