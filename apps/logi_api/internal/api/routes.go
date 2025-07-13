package api

import (
	"bombayv/logiapp-monorepo/logi_api/internal/api/handlers"
	"bombayv/logiapp-monorepo/logi_api/internal/api/middleware"
	"bombayv/logiapp-monorepo/logi_api/internal/config"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/cache"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/database"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with all the necessary routes and middleware.
func SetupRouter(db *database.DB, redisCache *cache.Cache, userHandler *handlers.UserHandler) *gin.Engine {
	gin.SetMode(config.App.GinMode)
	router := gin.Default()
	if config.App.TrustedProxy != "" {
		err := router.SetTrustedProxies([]string{config.App.TrustedProxy})
		if err != nil {
			log.Fatalf("could not set trusted proxies: %v", err)
			return nil
		}
	}

	router.Use(middleware.RequestLogger(5 * time.Second))
	router.Use(func(c *gin.Context) {
		c.Set("cache", redisCache)
		c.Next()
	})

	router.GET("/status", handlers.Status)

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Public routes

		// User routes (protected)
		v1.POST("/users/login", userHandler.Login)
		v1.POST("/users/logout", middleware.AuthMiddleware([]string{"sales", "driver", "admin"}), userHandler.Logout)
		v1.GET("/users/me", middleware.AuthMiddleware([]string{"sales", "driver", "admin"}), userHandler.Me)
		v1.POST("/users/register", middleware.AuthMiddleware([]string{"admin", "sales"}), userHandler.RegisterUser)

		// Test route for revoked tokens
		v1.GET("/test-revoked-token", middleware.AuthMiddleware([]string{}), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Access granted: Token is valid and not revoked."})
		})
	}

	return router
}
