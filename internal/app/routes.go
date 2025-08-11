package app

import (
	"haslaw-be-services/internal/middleware"
	"haslaw-be-services/internal/models"

	"github.com/gin-gonic/gin"
)

func (a *App) setupPublicRoutes(v1 *gin.RouterGroup) {
	authHandler := a.getAuthHandler()

	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

}

func (a *App) setupAuthRoutes(v1 *gin.RouterGroup) {
	authHandler := a.getAuthHandler()
	adminHandler := a.getAdminHandler()
	authService := a.getAuthService()

	auth := v1.Group("/auth")
	auth.Use(middleware.AuthMiddleware(authService))
	{
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/profile", adminHandler.GetProfile)
		auth.PUT("/profile", adminHandler.UpdateProfile)
	}
}

func (a *App) setupAdminRoutes(v1 *gin.RouterGroup) {
	authService := a.getAuthService()
	newsHandler := a.getNewsHandler()
	memberHandler := a.getMemberHandler()

	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(authService))
	admin.Use(middleware.RequireRole(models.Admin))
	{

		news := admin.Group("/news")
		{
			news.POST("", newsHandler.Create)
			news.PUT("/:id", newsHandler.Update)
			news.DELETE("/:id", newsHandler.Delete)
		}

		members := admin.Group("/members")
		{
			members.POST("", memberHandler.Create)
			members.PUT("/:id", memberHandler.Update)
			members.DELETE("/:id", memberHandler.Delete)
		}
	}
}

func (a *App) setupSuperAdminRoutes(v1 *gin.RouterGroup) {
	authService := a.getAuthService()
	adminHandler := a.getAdminHandler()

	superAdmin := v1.Group("/super-admin")
	superAdmin.Use(middleware.AuthMiddleware(authService))
	superAdmin.Use(middleware.RequireRole(models.SuperAdmin))
	{

		admins := superAdmin.Group("/admins")
		{
			admins.POST("", adminHandler.CreateAdmin)
		}
	}
}
