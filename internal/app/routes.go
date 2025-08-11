package app

import (
	"haslaw-be-services/internal/middleware"
	"haslaw-be-services/internal/models"

	"github.com/gin-gonic/gin"
)

// setupPublicRoutes sets up routes that don't require authentication
func (a *App) setupPublicRoutes(v1 *gin.RouterGroup) {
	authHandler := a.getAuthHandler()
	newsHandler := a.getNewsHandler()
	memberHandler := a.getMemberHandler()

	// Auth routes (public)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Public news routes
	news := v1.Group("/news")
	{
		news.GET("", newsHandler.GetAllPublicNews)     // Get all published news
		news.GET("/:id", newsHandler.GetByID)          // Get news by ID
		news.GET("/slug/:slug", newsHandler.GetBySlug) // Get news by slug
	}

	// Public member routes
	members := v1.Group("/members")
	{
		members.GET("", memberHandler.GetAll)
		members.GET("/:id", memberHandler.GetByID)
	}
}

// setupAuthRoutes sets up routes that require authentication (admin or super admin)
func (a *App) setupAuthRoutes(v1 *gin.RouterGroup) {
	authHandler := a.getAuthHandler()
	adminHandler := a.getAdminHandler()
	authService := a.getAuthService()

	// Routes requiring authentication
	auth := v1.Group("/auth")
	auth.Use(middleware.AuthMiddleware(authService))
	{
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/profile", adminHandler.GetProfile)
		auth.PUT("/profile", adminHandler.UpdateProfile)
	}
}

// setupAdminRoutes sets up routes that require admin role or higher
func (a *App) setupAdminRoutes(v1 *gin.RouterGroup) {
	authService := a.getAuthService()
	newsHandler := a.getNewsHandler()
	memberHandler := a.getMemberHandler()

	// Admin routes (admin and super admin can access)
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authService))
	admin.Use(middleware.RequireRole(models.Admin))
	{
		// News management - CRUD lengkap
		news := admin.Group("/news")
		{
			news.GET("", newsHandler.GetAll)                           // Get all news (admin view)
			news.GET("/:id", newsHandler.GetByID)                      // Get specific news
			news.POST("", newsHandler.Create)                          // Create news
			news.PUT("/:id", newsHandler.Update)                       // Update news
			news.DELETE("/:id", newsHandler.Delete)                    // Delete news
			news.GET("/drafts", newsHandler.GetDrafts)                 // Get draft news
			news.GET("/drafts/:id", newsHandler.GetDraftByID)          // Get draft by ID
			news.POST("/drafts/:id/publish", newsHandler.PublishDraft) // Publish draft
		}

		// Member management - CRUD lengkap
		members := admin.Group("/members")
		{
			members.GET("", memberHandler.GetAll)        // Get all members
			members.GET("/:id", memberHandler.GetByID)   // Get specific member
			members.POST("", memberHandler.Create)       // Create member
			members.PUT("/:id", memberHandler.Update)    // Update member
			members.DELETE("/:id", memberHandler.Delete) // Delete member
		}
	}
}

// setupSuperAdminRoutes sets up routes that require super admin role
func (a *App) setupSuperAdminRoutes(v1 *gin.RouterGroup) {
	authService := a.getAuthService()
	adminHandler := a.getAdminHandler()

	// Super admin routes (only super admin can access)
	superAdmin := v1.Group("/super-admin")
	superAdmin.Use(middleware.AuthMiddleware(authService))
	superAdmin.Use(middleware.RequireRole(models.SuperAdmin))
	{
		// Admin management - hanya yang tersedia
		admins := superAdmin.Group("/admins")
		{
			admins.POST("", adminHandler.CreateAdmin) // Create new admin
		}
	}
}
