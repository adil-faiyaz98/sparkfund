package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	rate   rate.Limit
	burst  int
	expiry time.Duration
}

func NewRateLimiter(r rate.Limit, b int, expiry time.Duration) *RateLimiter {
	return &RateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		rate:   r,
		burst:  b,
		expiry: expiry,
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.ips[ip] = limiter

		// Clean up old entries
		time.AfterFunc(rl.expiry, func() {
			rl.mu.Lock()
			delete(rl.ips, ip)
			rl.mu.Unlock()
		})
	}

	return limiter
}

func RateLimit(requests int, per time.Duration, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rate.Limit(float64(requests)/per.Seconds()), burst, per*2)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
				"rate": gin.H{
					"requests": requests,
					"per":      per.String(),
					"burst":    burst,
				},
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
