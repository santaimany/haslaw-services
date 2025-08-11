package main

import (
	"haslaw-be-services/internal/config"
	"haslaw-be-services/internal/models"
	"haslaw-be-services/internal/utils"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("🚀 Starting complete database migration...")

	// Step 1: Auto migrate all models
	log.Println("📋 Creating/updating all tables...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.News{},
		&models.Member{},
		&models.BlacklistedToken{},
	); err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}
	log.Println("✅ All tables created/updated successfully!")

	// Step 2: Ensure role column exists and has proper default
	log.Println("🔧 Ensuring role column configuration...")
	if err := db.Exec("ALTER TABLE users MODIFY COLUMN role VARCHAR(20) NOT NULL DEFAULT 'admin'").Error; err != nil {
		log.Printf("⚠️  Warning: Could not modify role column (might be correct already): %v", err)
	}

	// Step 3: Update existing users without role
	log.Println("👥 Updating existing users with default admin role...")
	result := db.Exec("UPDATE users SET role = 'admin' WHERE role IS NULL OR role = ''")
	if result.Error != nil {
		log.Printf("⚠️  Warning: Could not update existing users: %v", result.Error)
	} else {
		log.Printf("✅ Updated %d existing users with admin role", result.RowsAffected)
	}

	// Step 4: Create default super admin if not exists
	log.Println("👑 Creating default super admin...")
	var existingSuperAdmin models.User
	err = db.Where("username = ?", "superadmin").First(&existingSuperAdmin).Error

	if err != nil {
		// Super admin doesn't exist, create one
		hashedPassword, err := utils.HashPassword("superadmin123")
		if err != nil {
			log.Fatal("❌ Failed to hash super admin password:", err)
		}

		superAdmin := models.User{
			Username: "superadmin",
			Email:    "superadmin@haslaw.com",
			Password: hashedPassword,
			Role:     models.SuperAdmin,
		}

		if err := db.Create(&superAdmin).Error; err != nil {
			log.Fatal("❌ Failed to create super admin:", err)
		}

		log.Println("✅ Default super admin created successfully!")
		log.Println("📝 Super Admin Credentials:")
		log.Println("   Username: superadmin")
		log.Println("   Password: superadmin123")
		log.Println("   ⚠️  IMPORTANT: Please change this password after first login!")
	} else {
		log.Println("✅ Super admin already exists, skipping creation")
	}

	// Step 5: Create sample news if none exist (optional)
	log.Println("📰 Checking for sample news...")
	var newsCount int64
	db.Model(&models.News{}).Count(&newsCount)

	if newsCount == 0 {
		log.Println("📝 Creating sample news articles...")
		sampleNews := []models.News{
			{
				NewsTitle: "Welcome to HasLaw Services",
				Slug:      "welcome-to-haslaw-services",
				Category:  "Company News",
				Status:    models.Posted,
				Content:   "Welcome to our new content management system. This is a sample news article to demonstrate the functionality.",
			},
			{
				NewsTitle: "Legal Update - Corporate Law Changes",
				Slug:      "legal-update-corporate-law-changes",
				Category:  "Legal Updates",
				Status:    models.Drafted,
				Content:   "This is a draft article about recent changes in corporate law. It will be published after review.",
			},
		}

		for _, news := range sampleNews {
			if err := db.Create(&news).Error; err != nil {
				log.Printf("⚠️  Warning: Could not create sample news '%s': %v", news.NewsTitle, err)
			}
		}
		log.Printf("✅ Created %d sample news articles", len(sampleNews))
	} else {
		log.Printf("✅ Found %d existing news articles, skipping sample creation", newsCount)
	}

	// Step 6: Cleanup expired blacklisted tokens (if any)
	log.Println("🧹 Cleaning up expired blacklisted tokens...")
	result = db.Where("expires_at < NOW()").Delete(&models.BlacklistedToken{})
	if result.Error != nil {
		log.Printf("⚠️  Warning: Could not cleanup expired tokens: %v", result.Error)
	} else {
		log.Printf("✅ Cleaned up %d expired tokens", result.RowsAffected)
	}

	// Step 7: Verify database structure
	log.Println("🔍 Verifying database structure...")

	// Check if all tables exist
	tables := []string{"users", "news", "members", "blacklisted_tokens"}
	for _, table := range tables {
		var count int64
		if err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", table).Scan(&count).Error; err != nil {
			log.Printf("❌ Error checking table %s: %v", table, err)
		} else if count == 0 {
			log.Printf("❌ Table %s not found!", table)
		} else {
			log.Printf("✅ Table %s exists", table)
		}
	}

	log.Println("")
	log.Println("🎉 =====================================")
	log.Println("🎉 DATABASE MIGRATION COMPLETED!")
	log.Println("🎉 =====================================")
	log.Println("📊 Summary:")
	log.Println("   - All tables created/updated")
	log.Println("   - User roles configured")
	log.Println("   - Super admin ready")
	log.Println("   - Sample data available")
	log.Println("   - Token blacklist system ready")
	log.Println("")
	log.Println("🚀 Your application is ready to run!")
	log.Println("💡 Start with: go run ./cmd/api")
}
