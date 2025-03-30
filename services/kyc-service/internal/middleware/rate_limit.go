package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RedisURL      string
	RequestsPerIP int
	WindowSeconds int
	BlockDuration int
}

// DefaultRateLimitConfig returns default rate limiting configuration
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RedisURL:      "redis:6379",
		RequestsPerIP: 100,    // 100 requests
		WindowSeconds: 60,     // per minute
		BlockDuration: 300,    // block for 5 minutes if exceeded
	}
}

// RateLimitMiddleware provides rate limiting functionality
func RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisURL,
	})

	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()

		// Check if IP is blocked
		blocked, err := rdb.Exists(ctx, fmt.Sprintf("block:%s", ip)).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
			c.Abort()
			return
		}
		if blocked == 1 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "IP is blocked due to excessive requests"})
			c.Abort()
			return
		}

		// Get current request count
		key := fmt.Sprintf("rate:%s", ip)
		count, err := rdb.Get(ctx, key).Int()
		if err == redis.Nil {
			// First request in window
			err = rdb.Set(ctx, key, 1, time.Duration(config.WindowSeconds)*time.Second).Err()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
				c.Abort()
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
			c.Abort()
			return
		} else if count >= config.RequestsPerIP {
			// Block IP if limit exceeded
			err = rdb.Set(ctx, fmt.Sprintf("block:%s", ip), 1, time.Duration(config.BlockDuration)*time.Second).Err()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
				c.Abort()
				return
			}
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		} else {
			// Increment counter
			err = rdb.Incr(ctx, key).Err()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
} 