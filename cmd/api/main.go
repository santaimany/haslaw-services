package main

import (
	"fmt"
	"log"
	"os"

	"haslaw-be-services/internal/app"
)

func main() {
	// Initialize application
	application, err := app.New()
	if err != nil {
		log.Fatal("❌ Gagal menjalankan aplikasi:", err)
	}
	defer application.Close()

	// Get port from config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Display startup message
	fmt.Printf("🚀 Server sedang berjalan di port %s\n", port)
	fmt.Println("📋 Endpoint yang tersedia:")
	fmt.Println("   🔍 GET /health                          -> Health check")
	fmt.Println("")
	fmt.Println("   📝 Auth Endpoints:")
	fmt.Println("   - POST /api/v1/auth/login               -> Login")
	fmt.Println("   - POST /api/v1/auth/refresh             -> Refresh token")
	fmt.Println("   - POST /api/v1/auth/logout              -> Logout (perlu auth)")
	fmt.Println("   - GET /api/v1/auth/profile              -> Lihat profil (perlu auth)")
	fmt.Println("   - PUT /api/v1/auth/profile              -> Update profil (perlu auth)")
	fmt.Println("")
	fmt.Println("   📰 Public News Endpoints:")
	fmt.Println("   - GET /api/v1/news                      -> Lihat semua berita")
	fmt.Println("   - GET /api/v1/news/:id                  -> Lihat berita by ID")
	fmt.Println("   - GET /api/v1/news/slug/:slug           -> Lihat berita by slug")
	fmt.Println("")
	fmt.Println("   👥 Public Member Endpoints:")
	fmt.Println("   - GET /api/v1/members                   -> Lihat semua anggota")
	fmt.Println("   - GET /api/v1/members/:id               -> Lihat anggota by ID")
	fmt.Println("")
	fmt.Println("   🔒 Admin News Management (perlu role admin+):")
	fmt.Println("   - GET /api/v1/admin/news                -> Lihat semua berita (admin)")
	fmt.Println("   - GET /api/v1/admin/news/:id            -> Lihat berita by ID (admin)")
	fmt.Println("   - POST /api/v1/admin/news               -> Buat berita baru")
	fmt.Println("   - PUT /api/v1/admin/news/:id            -> Update berita")
	fmt.Println("   - DELETE /api/v1/admin/news/:id         -> Hapus berita")
	fmt.Println("   - GET /api/v1/admin/news/drafts         -> Lihat draft berita")
	fmt.Println("   - GET /api/v1/admin/news/drafts/:id     -> Lihat draft by ID")
	fmt.Println("   - POST /api/v1/admin/news/drafts/:id/publish -> Publish draft")
	fmt.Println("")
	fmt.Println("   🔒 Admin Member Management (perlu role admin+):")
	fmt.Println("   - GET /api/v1/admin/members             -> Lihat semua anggota (admin)")
	fmt.Println("   - GET /api/v1/admin/members/:id         -> Lihat anggota by ID (admin)")
	fmt.Println("   - POST /api/v1/admin/members            -> Buat anggota baru")
	fmt.Println("   - PUT /api/v1/admin/members/:id         -> Update anggota")
	fmt.Println("   - DELETE /api/v1/admin/members/:id      -> Hapus anggota")
	fmt.Println("")
	fmt.Println("   👑 Super Admin Management (perlu role super_admin):")
	fmt.Println("   - POST /api/v1/super-admin/admins       -> Buat admin baru")
	fmt.Println("")
	fmt.Println("📚 Default Super Admin:")
	fmt.Println("   Username: superadmin")
	fmt.Println("   Password: superadmin123")
	// Start server
	if err := application.Router.Run(":" + port); err != nil {
		log.Fatal("❌ Server gagal berjalan:", err)
	}
}
