package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs details about each incoming request.
func RequestLogger(slowThreshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Log details after the request is handled
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		logMessage := ""
		if latency > slowThreshold {
			logMessage = "SLOW REQUEST"
		}

		log.Printf("[GIN] %s %3d | %13v | %15s | %-7s %s",
			logMessage,
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
