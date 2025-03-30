package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redis *redis.Client
}

func NewRateLimiter(redisAddr string) (*RateLimiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RateLimiter{redis: client}, nil
}

func (rl *RateLimiter) RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rate_limit:%s:%s", c.ClientIP(), c.Request.URL.Path)
		ctx := context.Background()

		// Get current count
		count, err := rl.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limit check failed"})
			c.Abort()
			return
		}

		// If key doesn't exist, set it with expiration
		if err == redis.Nil {
			err = rl.redis.Set(ctx, key, 1, window).Err()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limit check failed"})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// Check if limit exceeded
		if count >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"retry_after": window.Seconds(),
			})
			c.Abort()
			return
		}

		// Increment counter
		err = rl.redis.Incr(ctx, key).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limit check failed"})
			c.Abort()
			return
		}

		c.Next()
	}
} 