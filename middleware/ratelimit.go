package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimiter(requestsPerSecond int, burstSize int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:     rate.Limit(requestsPerSecond),
		burst:    burstSize,
	}
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check locking pattern
		limiter, exists = rl.visitors[ip]
		if !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.visitors[ip] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

func (rl *RateLimiter) cleanupVisitors() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Simple cleanup - reset map if it gets too big
	if len(rl.visitors) > 1000 {
		rl.visitors = make(map[string]*rate.Limiter)
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	// Clean up old visitors periodically
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
		clientIP := c.ClientIP()
		limiter := rl.getVisitor(clientIP)

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

// DefaultRateLimiter returns 10 req/sec with burst of 20
func DefaultRateLimiter() *RateLimiter {
	return NewRateLimiter(10, 20)
}

// StrictRateLimiter returns 5 req/sec with burst of 10
func StrictRateLimiter() *RateLimiter {
	return NewRateLimiter(5, 10)
}
