package middleware

import (
	"haslaw-be-services/internal/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type client struct {
	requests []time.Time
}

var (
	clients = make(map[string]*client)
	mutex   = &sync.Mutex{}
)

func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		mutex.Lock()
		defer mutex.Unlock()

		now := time.Now()

		// Get or create client
		if clients[clientIP] == nil {
			clients[clientIP] = &client{requests: []time.Time{}}
		}

		client := clients[clientIP]

		// Remove requests older than 1 minute
		var validRequests []time.Time
		for _, reqTime := range client.requests {
			if now.Sub(reqTime) < time.Minute {
				validRequests = append(validRequests, reqTime)
			}
		}
		client.requests = validRequests

		// Check if rate limit exceeded
		if len(client.requests) >= requestsPerMinute {
			utils.ErrorResponse(c, 429, "Rate limit exceeded", "Too many requests. Please try again later.")
			c.Abort()
			return
		}

		// Add current request
		client.requests = append(client.requests, now)

		c.Next()
	}
}
