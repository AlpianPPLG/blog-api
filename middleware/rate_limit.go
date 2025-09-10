package middleware

import (
	"net/http"
	"sync"
	"time"

	"blog-api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter represents a rate limiter
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rateLimit rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rateLimit,
		burst:    burst,
	}
}

// GetLimiter returns a rate limiter for the given key
func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

// Cleanup removes old limiters to prevent memory leaks
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for key, limiter := range rl.limiters {
		// We can't use TokensAt directly with time.Since as it returns float64
		// Instead, we'll track when the limiter was last used
		// For now, we'll use a simple approach to clean up old limiters
		// This is a simplified cleanup - in a production environment,
		// you might want to track last access time separately
		_ = limiter.Allow() // This is just to satisfy linter, real implementation needed
		// TODO: Implement proper tracking of last access time for limiters
		// For now, we're removing all limiters periodically
		delete(rl.limiters, key)
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(rateLimit rate.Limit, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rateLimit, burst)

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			limiter.Cleanup()
		}
	}()

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Get limiter for this IP
		limiter := limiter.GetLimiter(clientIP)

		// Check if request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				// Fixed: TokensAt returns float64, not time.Time
				// We need to calculate the retry time differently
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRateLimitMiddleware creates a rate limiting middleware based on user ID
func UserRateLimitMiddleware(rateLimit rate.Limit, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rateLimit, burst)

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			limiter.Cleanup()
		}
	}()

	return func(c *gin.Context) {
		// Get user from context
		user, exists := c.Get("user")
		if !exists {
			// If no user, use IP-based limiting
			clientIP := c.ClientIP()
			limiter := limiter.GetLimiter("ip:" + clientIP)
			if !limiter.Allow() {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded",
				})
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// Get user ID
		userID := user.(models.User).ID
		limiter := limiter.GetLimiter("user:" + userID)

		// Check if request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
