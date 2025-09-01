package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const requestIDKey = "req_id"

// RequestID üretir veya gelen X-Request-ID’yi kullanır.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = randHex(16)
		}
		c.Set(requestIDKey, rid)
		c.Writer.Header().Set("X-Request-ID", rid)
		c.Next()
	}
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// fallback basit (çok düşük olasılık)
		return "00000000000000000000000000000000"
	}
	return hex.EncodeToString(b)
}

// GetRequestID yardımcı
func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(requestIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
