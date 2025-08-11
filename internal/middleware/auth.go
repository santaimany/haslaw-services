package middleware

import (
	"haslaw-be-services/internal/models"
	"haslaw-be-services/internal/service"
	"haslaw-be-services/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := tokenParts[1]

		isBlacklisted, err := authService.IsTokenBlacklisted(token)
		if err != nil {
			utils.UnauthorizedResponse(c, "Token validation failed")
			c.Abort()
			return
		}

		if isBlacklisted {
			utils.UnauthorizedResponse(c, "Token has been invalidated")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("token", token)
		c.Next()
	}
}

func RequireRole(requiredRole models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.ForbiddenResponse(c, "User role not found")
			c.Abort()
			return
		}

		userRole := models.UserRole(role.(string))

		// Check role hierarchy: SuperAdmin can access Admin routes, but not vice versa
		hasPermission := false

		switch requiredRole {
		case models.SuperAdmin:
			// Only SuperAdmin can access SuperAdmin routes
			hasPermission = userRole == models.SuperAdmin
		case models.Admin:
			// Both Admin and SuperAdmin can access Admin routes
			hasPermission = userRole == models.Admin || userRole == models.SuperAdmin
		default:
			hasPermission = userRole == requiredRole
		}

		if !hasPermission {
			utils.ForbiddenResponse(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireSuperAdmin() gin.HandlerFunc {
	return RequireRole(models.SuperAdmin)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}
