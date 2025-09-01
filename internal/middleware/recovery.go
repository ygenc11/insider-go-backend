package middleware

import (
	"net/http"
	"runtime/debug"

	"log/slog"

	"github.com/gin-gonic/gin"
)

// Recovery, panik durumlarında 500 JSON döner ve hatayı loglar.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("http_panic", "error", r, "stack", string(debug.Stack()))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}
