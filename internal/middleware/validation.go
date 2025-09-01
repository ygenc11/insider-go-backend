package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ContentTypeJSON, application/json zorunlu kılar
func ContentTypeJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		ct := c.GetHeader("Content-Type")
		if ct == "" || (ct != "application/json" && ct != "application/json; charset=utf-8") {
			c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"error": "content-type must be application/json"})
			return
		}
		c.Next()
	}
}

// MaxBodyBytes, istek gövdesi için üst sınır belirler (bayt cinsinden)
func MaxBodyBytes(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Writer.Header().Set("X-Max-Body-Bytes", strconv.FormatInt(limit, 10))
		c.Next()
	}
}
