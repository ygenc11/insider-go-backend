package middleware

import (
	"net/http"
	"strings"

	"insider-go-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware, JWT token kontrolü yapar
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization header al
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// "Bearer <token>" formatını kontrol et
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// Token parse et
		userID, role, err := services.ParseJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Context’e ekle
		c.Set("user_id", userID)
		c.Set("role", role)

		c.Next()
	}
}

// RequireRole en az bir rolün eşleşmesini zorunlu kılar (ör: "admin").
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role := c.GetString("role")
		if _, ok := allowed[role]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
			c.Abort()
			return
		}
		c.Next()
	}
}
