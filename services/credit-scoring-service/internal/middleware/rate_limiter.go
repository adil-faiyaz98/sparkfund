package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sparkfund/credit-scoring-service/internal/errors"
	"go.uber.org/zap"
)

type RateLimiterConfig struct {
	RequestsPerMinute int
	BurstSize        int
	RedisURL         string
	RedisPassword    string
	RedisDB          int
}

type RateLimiter struct {
	config *RateLimiterConfig
	client *redis.Client
	logger *zap.Logger
}

func NewRateLimiter(config *RateLimiterConfig, logger *zap.Logger) (*RateLimiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisURL,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RateLimiter{
		config: config,
		client: client,
		logger: logger,
	}, nil
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (IP or user ID)
		clientID := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			clientID = fmt.Sprintf("user:%s", userID)
		}

		// Create Redis key
		key := fmt.Sprintf("rate_limit:%s", clientID)

		// Use Redis transaction for atomic operations
		err := rl.client.Watch(context.Background(), func(tx *redis.Tx) error {
			// Get current count
			count, err := tx.Get(context.Background(), key).Int()
			if err != nil && err != redis.Nil {
				return fmt.Errorf("failed to get rate limit count: %w", err)
			}

			// Check if rate limit exceeded
			if count >= rl.config.RequestsPerMinute {
				rl.logger.Warn("rate limit exceeded",
					zap.String("client_id", clientID),
					zap.Int("count", count),
				)
				c.AbortWithStatusJSON(http.StatusTooManyRequests, errors.NewAPIError(
					errors.ErrRateLimitExceeded,
					"rate limit exceeded",
				))
				return nil
			}

			// Increment counter
			pipe := tx.Pipeline()
			pipe.Incr(context.Background(), key)
			pipe.Expire(context.Background(), key, time.Minute)
			_, err = pipe.Exec(context.Background())
			if err != nil {
				return fmt.Errorf("failed to increment rate limit: %w", err)
			}

			return nil
		}, key).Err()

		if err != nil {
			rl.logger.Error("rate limit error", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, errors.NewAPIError(
				errors.ErrInternalServer,
				"rate limit error",
			))
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) Close() error {
	return rl.client.Close()
} 