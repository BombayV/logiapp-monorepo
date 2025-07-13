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

	router.GET("/status", handlers.Status)
	router.GET("/token", handlers.GetToken)

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Public routes

		// User routes (protected)
		v1.POST("/users/login", handlers.LoginUser)
		// Apply auth middleware to registration route
		v1.POST("/users/register", middleware.AuthMiddleware([]string{""}), userHandler.RegisterUser)
	}

	return router
}
