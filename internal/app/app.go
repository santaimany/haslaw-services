package app

import (
	"fmt"
	"os"

	"haslaw-be-services/internal/config"
	"haslaw-be-services/internal/handlers"
	"haslaw-be-services/internal/middleware"
	"haslaw-be-services/internal/models"
	"haslaw-be-services/internal/repository"
	"haslaw-be-services/internal/service"
	"haslaw-be-services/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func New() (*App, error) {

	if err := godotenv.Load(); err != nil {

		fmt.Println("⚠️  No .env file found, using system environment variables")
	}

	db, err := config.NewDatabase()
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	if err := utils.InitUploadDirectories(); err != nil {
		return nil, fmt.Errorf("upload directories initialization failed: %w", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := &App{
		DB:     db,
		Router: gin.New(),
	}

	if err := app.initializeServices(); err != nil {
		return nil, fmt.Errorf("service initialization failed: %w", err)
	}

	if err := app.setupMiddleware(); err != nil {
		return nil, fmt.Errorf("middleware setup failed: %w", err)
	}

	if err := app.setupRoutes(); err != nil {
		return nil, fmt.Errorf("routes setup failed: %w", err)
	}

	return app, nil
}

func (a *App) Close() error {

	sqlDB, err := a.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.News{},
		&models.Member{},
		&models.BlacklistedToken{},
	)
}

func (a *App) initializeServices() error {

	userRepo := repository.NewUserRepository(a.DB)
	blacklistRepo := repository.NewBlacklistRepository(a.DB)
	newsRepo := repository.NewNewsRepository(a.DB)
	memberRepo := repository.NewMemberRepository(a.DB)

	authService := service.NewAuthService(userRepo, blacklistRepo)
	newsService := service.NewNewsService(newsRepo)
	memberService := service.NewMemberService(memberRepo)

	if err := authService.CreateDefaultSuperAdmin(); err != nil {
		return fmt.Errorf("failed to create default super admin: %w", err)
	}

	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(authService, userRepo)
	newsHandler := handlers.NewNewsHandler(newsService)
	memberHandler := handlers.NewMemberHandler(memberService)
	healthHandler := handlers.NewHealthHandler()

	a.Router.Use(func(c *gin.Context) {
		c.Set("authHandler", authHandler)
		c.Set("adminHandler", adminHandler)
		c.Set("newsHandler", newsHandler)
		c.Set("memberHandler", memberHandler)
		c.Set("healthHandler", healthHandler)
		c.Set("authService", authService)
		c.Next()
	})

	return nil
}

func (a *App) setupMiddleware() error {
	// Performance optimizations
	a.Router.Use(gin.Recovery())

	// Only use logger in development
	if os.Getenv("GIN_MODE") != "release" {
		a.Router.Use(gin.Logger())
	}

	// Add compression for better response times
	a.Router.Use(middleware.GzipMiddleware())
	
	// Core middlewares
	a.Router.Use(middleware.TraceIDMiddleware())
	a.Router.Use(middleware.RateLimitMiddleware(100)) // Increase rate limit
	a.Router.Use(middleware.CORSMiddleware())
	a.Router.Use(middleware.SecurityHeadersMiddleware())

	// Static files with cache headers for better performance
	a.Router.Static("/uploads", "./uploads")

	return nil
}

func (a *App) setupRoutes() error {

	healthHandler := a.getHealthHandler()
	a.Router.GET("/health", healthHandler.Check)

	v1 := a.Router.Group("/api/v1")

	a.setupPublicRoutes(v1)
	a.setupAuthRoutes(v1)
	a.setupAdminRoutes(v1)
	a.setupSuperAdminRoutes(v1)

	return nil
}

func (a *App) getAuthHandler() *handlers.AuthHandler {
	userRepo := repository.NewUserRepository(a.DB)
	blacklistRepo := repository.NewBlacklistRepository(a.DB)
	authService := service.NewAuthService(userRepo, blacklistRepo)
	return handlers.NewAuthHandler(authService)
}

func (a *App) getAdminHandler() *handlers.AdminHandler {
	userRepo := repository.NewUserRepository(a.DB)
	blacklistRepo := repository.NewBlacklistRepository(a.DB)
	authService := service.NewAuthService(userRepo, blacklistRepo)
	return handlers.NewAdminHandler(authService, userRepo)
}

func (a *App) getNewsHandler() *handlers.NewsHandler {
	newsRepo := repository.NewNewsRepository(a.DB)
	newsService := service.NewNewsService(newsRepo)
	return handlers.NewNewsHandler(newsService)
}

func (a *App) getMemberHandler() *handlers.MemberHandler {
	memberRepo := repository.NewMemberRepository(a.DB)
	memberService := service.NewMemberService(memberRepo)
	return handlers.NewMemberHandler(memberService)
}

func (a *App) getHealthHandler() *handlers.HealthHandler {
	return handlers.NewHealthHandler()
}

func (a *App) getAuthService() service.AuthService {
	userRepo := repository.NewUserRepository(a.DB)
	blacklistRepo := repository.NewBlacklistRepository(a.DB)
	return service.NewAuthService(userRepo, blacklistRepo)
}
