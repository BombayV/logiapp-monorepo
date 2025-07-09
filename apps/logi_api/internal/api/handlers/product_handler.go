package handlers

import (
	"bombayv/logiapp-monorepo/logi_api/internal/storage/cache"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProducts handles the request to fetch all products.
func GetProducts(db *database.DB, redisCache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a real app, you would call the product service to get data.
		// For now, we'll return some mock data.
		mockProducts := []gin.H{
			{"id": "prod_1", "name": "Super Widget", "price": 19.99},
			{"id": "prod_2", "name": "Mega Gadget", "price": 29.99},
		}
		c.JSON(http.StatusOK, mockProducts)
	}
}

// GetProductByID handles the request to fetch a single product by its ID.
func GetProductByID(db *database.DB, redisCache *cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")

		// In a real app, you'd fetch this from the product service.
		// We'll just return a mock response.
		c.JSON(http.StatusOK, gin.H{
			"id":    productID,
			"name":  "Sample Product",
			"price": 49.99,
		})
	}
}
