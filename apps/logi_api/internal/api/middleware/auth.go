package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides a simple authentication check.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a real app, you would parse and validate a JWT.
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check for "Bearer " prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		token := parts[1]
		// Here you would validate the token (e.g., using a JWT library)
		if token != "fake-jwt-token-for-testing" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set user info in the context for downstream handlers
		c.Set("userID", "user_123")

		c.Next()
	}
}
