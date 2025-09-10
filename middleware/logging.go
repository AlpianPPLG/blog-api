package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware creates a logging middleware
func LoggingMiddleware(logger interface{}) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Printf("HTTP Request - Method: %s, Path: %s, Status: %d, Latency: %v, ClientIP: %s",
			param.Method, param.Path, param.StatusCode, param.Latency, param.ClientIP)
		return ""
	})
}

// RequestLoggingMiddleware creates a custom request logging middleware
func RequestLoggingMiddleware(logger interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("HTTP Request - Method: %s, Path: %s, Status: %d, Latency: %v, ClientIP: %s, BodySize: %d",
			method, path, statusCode, latency, clientIP, bodySize)
	}
}
