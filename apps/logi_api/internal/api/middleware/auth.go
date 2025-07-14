package middleware

import (
	"bombayv/logiapp-monorepo/logi_api/internal/config"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/cache"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens.
func AuthMiddleware(role []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: no token provided"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		tokenString = tokenString[len("Bearer "):]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.App.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// Check if the token is revoked
			cacheVal, exists := c.Get("cache")
			if !exists {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}
			isRevoked, err := cacheVal.(*cache.Cache).IsTokenRevoked(c.Request.Context(), claims.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}
			if isRevoked {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
				c.Abort()
				return
			}

			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)
			// Check if the user is admin or if the role is allowed
			if !slices.Contains(role, claims.Role) && claims.Role != "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
				c.Abort()
				return
			}

			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}
