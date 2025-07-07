package api

import (
	"entityevolution/ev-logiapp-monorepo/logi_api/internal/api/handlers"
	"entityevolution/ev-logiapp-monorepo/logi_api/internal/api/middleware"
	"entityevolution/ev-logiapp-monorepo/logi_api/internal/storage/cache"
	"entityevolution/ev-logiapp-monorepo/logi_api/internal/storage/database"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with all the necessary routes and middleware.
func SetupRouter(db *database.DB, redisCache *cache.Cache) *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	// Add custom logger middleware
	router.Use(middleware.RequestLogger(5 * time.Second))

	// Health check route
	router.GET("/status", handlers.Status)

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// User routes
		userRoutes := v1.Group("/users")
		{
			// Public routes
			userRoutes.POST("/register", handlers.RegisterUser(db))
			userRoutes.POST("/login", handlers.LoginUser(db))
		}

		// Product routes
		productRoutes := v1.Group("/products")
		// Apply auth middleware to all product routes
		productRoutes.Use(middleware.AuthMiddleware())
		{
			productRoutes.GET("", handlers.GetProducts(db, redisCache))
			productRoutes.GET("/:id", handlers.GetProductByID(db, redisCache))
		}
	}

	return router
}
