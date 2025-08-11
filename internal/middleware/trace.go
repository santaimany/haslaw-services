package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// TraceIDKey is the key used to store trace ID in context
const TraceIDKey = "trace_id"

// GenerateTraceID creates a unique trace ID
func GenerateTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if crypto/rand fails
		return hex.EncodeToString([]byte(time.Now().Format("20060102150405.000000")))
	}
	return hex.EncodeToString(bytes)
}

// TraceIDMiddleware adds a unique trace ID to each request
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or get trace ID from header
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = GenerateTraceID()
		}

		// Store trace ID in context
		c.Set(TraceIDKey, traceID)

		// Add trace ID to response header
		c.Header("X-Trace-ID", traceID)

		// Log request with trace ID
		start := time.Now()
		log.Printf("[TRACE: %s] %s %s - Started", traceID, c.Request.Method, c.Request.URL.Path)

		// Process request
		c.Next()

		// Log response with trace ID
		duration := time.Since(start)
		status := c.Writer.Status()
		log.Printf("[TRACE: %s] %s %s - Completed in %v with status %d",
			traceID, c.Request.Method, c.Request.URL.Path, duration, status)
	}
}

// GetTraceID extracts trace ID from gin context
func GetTraceID(c *gin.Context) string {
	if traceID, exists := c.Get(TraceIDKey); exists {
		return traceID.(string)
	}
	return "unknown"
}
