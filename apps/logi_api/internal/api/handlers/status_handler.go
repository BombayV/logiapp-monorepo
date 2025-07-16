package handlers

import (
	"context"
	"net/http"
	"time"

	"bombayv/logiapp-monorepo/logi_api/internal/storage/cache"
	"bombayv/logiapp-monorepo/logi_api/internal/storage/database"

	"github.com/gin-gonic/gin"
)

// StatusHandler handles status-related operations
type StatusHandler struct {
	DB    *database.DB
	Cache *cache.Cache
}

// NewStatusHandler creates a new status handler
func NewStatusHandler(db *database.DB, cache *cache.Cache) *StatusHandler {
	return &StatusHandler{
		DB:    db,
		Cache: cache,
	}
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	Name    string        `json:"name"`
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Latency time.Duration `json:"latency,omitempty"`
	Details interface{}   `json:"details,omitempty"`
}

// HealthResponse represents the overall health response
type HealthResponse struct {
	Status    string          `json:"status"`
	Timestamp time.Time       `json:"timestamp"`
	Services  []ServiceStatus `json:"services"`
}

// Health provides comprehensive health check for all services
func (h *StatusHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	services := []ServiceStatus{}
	overallStatus := "healthy"

	// Check database connectivity
	dbStatus := h.checkDatabase(ctx)
	services = append(services, dbStatus)
	if dbStatus.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// Check cache connectivity
	cacheStatus := h.checkCache(ctx)
	services = append(services, cacheStatus)
	if cacheStatus.Status != "healthy" {
		overallStatus = "degraded"
	}

	// Check application health
	appStatus := h.checkApplication(ctx)
	services = append(services, appStatus)

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Services:  services,
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	} else if overallStatus == "degraded" {
		statusCode = http.StatusPartialContent
	}

	c.JSON(statusCode, response)
}

// Simple status endpoint for basic health checks
func (h *StatusHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now(),
	})
}

// checkDatabase verifies database connectivity and performance
func (h *StatusHandler) checkDatabase(ctx context.Context) ServiceStatus {
	start := time.Now()

	// Test basic connectivity
	err := h.DB.Pool.Ping(ctx)
	if err != nil {
		return ServiceStatus{
			Name:    "database",
			Status:  "unhealthy",
			Message: "Database connection failed",
			Latency: time.Since(start),
		}
	}

	// Test basic query functionality with a simple query
	var exists bool
	err = h.DB.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users LIMIT 1)").Scan(&exists)
	if err != nil {
		return ServiceStatus{
			Name:    "database",
			Status:  "unhealthy",
			Message: "Database query failed",
			Latency: time.Since(start),
		}
	}

	return ServiceStatus{
		Name:    "database",
		Status:  "healthy",
		Message: "Database is operational",
		Latency: time.Since(start),
	}
}

// checkCache verifies cache connectivity and performance
func (h *StatusHandler) checkCache(ctx context.Context) ServiceStatus {
	start := time.Now()

	// Test basic connectivity with a ping
	err := h.Cache.Client.Do(ctx, h.Cache.Client.B().Ping().Build()).Error()
	if err != nil {
		return ServiceStatus{
			Name:    "cache",
			Status:  "unhealthy",
			Message: "Cache connection failed",
			Latency: time.Since(start),
		}
	}

	// Test basic set/get functionality
	testKey := "health_check_test"
	testValue := "ok"

	// Set a test value
	err = h.Cache.Client.Do(ctx, h.Cache.Client.B().Set().Key(testKey).Value(testValue).Ex(time.Second*10).Build()).Error()
	if err != nil {
		return ServiceStatus{
			Name:    "cache",
			Status:  "unhealthy",
			Message: "Cache write operation failed",
			Latency: time.Since(start),
		}
	}

	// Get the test value
	result, err := h.Cache.Client.Do(ctx, h.Cache.Client.B().Get().Key(testKey).Build()).ToString()
	if err != nil || result != testValue {
		return ServiceStatus{
			Name:    "cache",
			Status:  "unhealthy",
			Message: "Cache read operation failed",
			Latency: time.Since(start),
		}
	}

	// Clean up test key
	h.Cache.Client.Do(ctx, h.Cache.Client.B().Del().Key(testKey).Build())

	return ServiceStatus{
		Name:    "cache",
		Status:  "healthy",
		Message: "Cache is operational",
		Latency: time.Since(start),
	}
}

// checkApplication verifies application-specific health
func (h *StatusHandler) checkApplication(ctx context.Context) ServiceStatus {
	return ServiceStatus{
		Name:    "application",
		Status:  "healthy",
		Message: "Application is running",
	}
}
