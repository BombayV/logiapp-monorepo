package main

import (
	"bombayv/logiapp-monorepo/logi_api/internal/api"
	"bombayv/logiapp-monorepo/logi_api/internal/api/handlers"
	"bombayv/logiapp-monorepo/logi_api/internal/config"
	"bombayv/logiapp-monorepo/logi_api/internal/core/orders"
	"bombayv/logiapp-monorepo/logi_api/internal/core/user"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/cache"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/database"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/repository"
	"fmt"
	"log"
)

func main() {
	// Load application configuration
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Initialize database, cache, etc. (dependencies)
	db, err := database.NewDatabase(config.App.DBSource)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	redisCache, err := cache.NewCache(config.App.RedisAddress)
	if err != nil {
		log.Fatalf("could not connect to cache: %v", err)
	}
	defer redisCache.Close()

	// Initialize repositories
	repoManager := repository.NewManager(db)

	// Initialize services
	userService := user.NewService(repoManager.User, redisCache)
	orderService := orders.NewService(repoManager.Orders, repoManager.User)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService)

	fmt.Println("Starting server on port", config.App.ServerPort)

	// Setup router
	router := api.SetupRouter(db, redisCache, userHandler, orderHandler)

	// Start server
	err = router.Run(":" + config.App.ServerPort)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
