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
func SetupRouter(db *database.DB, redisCache *cache.Cache, userHandler *handlers.UserHandler, orderHandler *handlers.OrderHandler) *gin.Engine {
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

		// User routes
		v1.GET("/users/me", middleware.AuthMiddleware([]string{"sales", "driver"}), userHandler.Me)
		v1.POST("/users/logout", middleware.AuthMiddleware([]string{"sales", "driver"}), userHandler.Logout)
		v1.POST("/users/login", userHandler.Login)

		// Get all users (admin only)
		v1.GET("/users", middleware.AuthMiddleware([]string{"sales"}), userHandler.GetUsers)
		v1.POST("/users", middleware.AuthMiddleware([]string{}), userHandler.RegisterUser)
		v1.GET("/users/:id", middleware.AuthMiddleware([]string{}), userHandler.GetUserByID)
		v1.DELETE("/users/:id", middleware.AuthMiddleware([]string{}), userHandler.DeleteUser)

		// Order routes (protected)
		v1.POST("/orders", middleware.AuthMiddleware([]string{"sales"}), orderHandler.CreateOrder)
		v1.GET("/orders/:id", middleware.AuthMiddleware([]string{"sales"}), orderHandler.GetOrderByID)
		v1.PUT("/orders/:id", middleware.AuthMiddleware([]string{"sales"}), orderHandler.UpdateOrder)
		v1.DELETE("/orders/:id", middleware.AuthMiddleware([]string{"sales"}), orderHandler.DeleteOrder)

		// Order items routes
		v1.POST("/orders/:order_id/items", middleware.AuthMiddleware([]string{"sales"}), orderHandler.AddOrderItem)
		v1.POST("/orders/:order_id/items/bulk", middleware.AuthMiddleware([]string{"sales"}), orderHandler.AddOrderItems)
		v1.GET("/orders/:order_id/items", middleware.AuthMiddleware([]string{"sales"}), orderHandler.GetOrderItems)
		v1.PUT("/orders/:order_id/items/:item_id", middleware.AuthMiddleware([]string{"sales"}), orderHandler.UpdateOrderItem)
		v1.DELETE("/orders/:order_id/items/:item_id", middleware.AuthMiddleware([]string{"sales"}), orderHandler.DeleteOrderItem)

		orders := v1.Group("/orders")
		{
			orders.GET("/", middleware.AuthMiddleware([]string{"sales", "admin"}), orderHandler.GetAllOrders)
			orders.POST("/", orderHandler.CreateOrder)
		}
	}

	return router
}
