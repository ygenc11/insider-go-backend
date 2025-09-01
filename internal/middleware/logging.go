package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger her HTTP isteğini slog kullanarak loglar.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)
		slog.Info("http_request",
			"method", method,
			"path", path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"client_ip", clientIP,
		)
	}
}
