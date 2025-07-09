package main

import (
	"bombayv/logiapp-monorepo/logi_api/internal/api"
	"bombayv/logiapp-monorepo/logi_api/internal/config"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/cache"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/database"
	"fmt"
	"log"
)

func main() {
	// Load application configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Initialize database, cache, etc. (dependencies)
	db, err := database.NewDatabase(cfg.DBSource)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	redisCache, err := cache.NewCache(cfg.RedisAddress)
	if err != nil {
		log.Fatalf("could not connect to cache: %v", err)
	}
	defer redisCache.Close()

	fmt.Println("Starting server on port", cfg.ServerPort)

	// Setup router
	router := api.SetupRouter(db, redisCache, cfg)

	// Start server
	err = router.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
