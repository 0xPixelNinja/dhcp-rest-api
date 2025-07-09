package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter holds the rate limiter configuration
type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter with specified rate and burst
func NewRateLimiter(requestsPerSecond int, burstSize int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burstSize,
	}
}

// getVisitor retrieves or creates a rate limiter for a specific IP
func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check after acquiring write lock
		limiter, exists = rl.visitors[ip]
		if !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.visitors[ip] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

// cleanupVisitors removes old visitors to prevent memory leaks
func (rl *RateLimiter) cleanupVisitors() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Create a new map to avoid memory leaks
	// In production, you might want to implement a more sophisticated cleanup
	// based on last access time
	if len(rl.visitors) > 1000 {
		rl.visitors = make(map[string]*rate.Limiter)
	}
}

// Middleware returns a Gin middleware function for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				rl.cleanupVisitors()
			}
		}
	}()

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Get rate limiter for this IP
		limiter := rl.getVisitor(clientIP)

		// Check if request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// DefaultRateLimiter creates a rate limiter with sensible defaults for production
// 10 requests per second with a burst of 20
func DefaultRateLimiter() *RateLimiter {
	return NewRateLimiter(10, 20)
}

// StrictRateLimiter creates a more restrictive rate limiter
// 5 requests per second with a burst of 10
func StrictRateLimiter() *RateLimiter {
	return NewRateLimiter(5, 10)
}
