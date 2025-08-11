package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GzipMiddleware provides gzip compression
func GzipMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if !shouldCompress(c.Request) {
			c.Next()
			return
		}

		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		gzipWriter := &gzipResponseWriter{
			ResponseWriter: c.Writer,
			Writer:        gz,
		}
		c.Writer = gzipWriter

		c.Next()
	})
}

type gzipResponseWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}

func shouldCompress(req *http.Request) bool {
	// Check if client accepts gzip encoding
	return strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")
}

// CacheMiddleware adds cache headers for static content
func CacheMiddleware(maxAge int) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Set cache headers for better performance
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		c.Header("ETag", generateETag(c.Request.URL.Path))
		
		// Check if client has cached version
		if match := c.GetHeader("If-None-Match"); match != "" {
			etag := c.GetHeader("ETag")
			if match == etag {
				c.AbortWithStatus(http.StatusNotModified)
				return
			}
		}
		
		c.Next()
	})
}

func generateETag(path string) string {
	// Simple ETag generation based on path and timestamp
	return fmt.Sprintf(`"%x"`, hash(path))
}

func hash(s string) uint32 {
	h := uint32(0)
	for _, c := range s {
		h = h*31 + uint32(c)
	}
	return h
}
