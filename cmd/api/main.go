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
	fmt.Println("   - GET /health                        -> Health check")
	fmt.Println("   - POST /api/v1/auth/login           -> Login")
	fmt.Println("   - POST /api/v1/auth/logout          -> Logout (perlu auth)")
	fmt.Println("   - GET /api/v1/auth/profile          -> Lihat profil (perlu auth)")
	fmt.Println("   - PUT /api/v1/auth/profile          -> Update profil (perlu auth)")
	fmt.Println("   - GET /api/v1/news                  -> Lihat semua berita")
	fmt.Println("   - GET /api/v1/members               -> Lihat semua anggota")
	fmt.Println("   - POST /api/v1/admin/news           -> Buat berita (admin+)")
	fmt.Println("   - POST /api/v1/admin/members        -> Buat anggota (admin+)")
	fmt.Println("   - POST /api/v1/super-admin/admins   -> Buat admin baru (super admin)")
	fmt.Println("")
	fmt.Println("📚 Default Super Admin:")
	fmt.Println("   Email: superadmin@haslaw.com")
	fmt.Println("   Password: superadmin123")

	// Start server
	if err := application.Router.Run(":" + port); err != nil {
		log.Fatal("❌ Server gagal berjalan:", err)
	}
}
