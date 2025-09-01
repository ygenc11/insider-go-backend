package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// basit token bucket
type bucket struct {
	tokens     float64
	lastRefill time.Time
}

// RateLimiter konfig
type RateLimiterConfig struct {
	RefillRatePerSec float64 // saniyede dolan token
	Burst            float64 // maksimum token
}

// RateLimiter IP başına kova uygular
func RateLimiter(cfg RateLimiterConfig) gin.HandlerFunc {
	var (
		mu sync.Mutex
		m  = map[string]*bucket{}
	)
	if cfg.RefillRatePerSec <= 0 {
		cfg.RefillRatePerSec = 10
	}
	if cfg.Burst <= 0 {
		cfg.Burst = 20
	}
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		mu.Lock()
		b, ok := m[ip]
		if !ok {
			b = &bucket{tokens: cfg.Burst, lastRefill: now}
			m[ip] = b
		}
		// refill
		dt := now.Sub(b.lastRefill).Seconds()
		b.tokens = minFloat(cfg.Burst, b.tokens+dt*cfg.RefillRatePerSec)
		b.lastRefill = now
		// consume
		if b.tokens < 1 {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		b.tokens -= 1
		mu.Unlock()
		c.Next()
	}
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
