package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"haslaw-be-services/internal/config"
)

// Temporary struct to work with old string format
type MemberOld struct {
	ID        uint   `gorm:"primaryKey"`
	Education string `gorm:"column:education"`
}

// Temporary struct to work with new array format
type MemberNew struct {
	ID        uint     `gorm:"primaryKey"`
	Education []string `gorm:"column:education;serializer:json"`
}

func main() {
	fmt.Println("🔄 Starting Education Migration...")

	// Connect to database
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	fmt.Println("✅ Database connected successfully!")

	// Get all members with old string education format
	var members []MemberOld
	if err := db.Table("members").Select("id, education").Find(&members).Error; err != nil {
		log.Fatal("❌ Failed to fetch members:", err)
	}

	fmt.Printf("📊 Found %d members to migrate\n", len(members))

	// Convert each member's education from string to array
	for _, member := range members {
		if member.Education == "" {
			continue // Skip empty education
		}

		fmt.Printf("🔄 Converting member ID %d...\n", member.ID)

		// Convert string to array (split by common separators)
		var educationArray []string

		// Try to parse as JSON first (in case some are already converted)
		var tempArray []string
		if err := json.Unmarshal([]byte(member.Education), &tempArray); err == nil {
			educationArray = tempArray
			fmt.Printf("   ✅ Already in JSON format\n")
		} else {
			// Split by common separators
			education := strings.ReplaceAll(member.Education, ", ", ",")
			education = strings.ReplaceAll(education, "; ", ";")

			if strings.Contains(education, ",") {
				educationArray = strings.Split(education, ",")
			} else if strings.Contains(education, ";") {
				educationArray = strings.Split(education, ";")
			} else {
				// Single education entry
				educationArray = []string{education}
			}

			// Clean up each entry
			for i, edu := range educationArray {
				educationArray[i] = strings.TrimSpace(edu)
			}

			fmt.Printf("   ✅ Converted to %d entries\n", len(educationArray))
		}

		// Convert to JSON
		educationJSON, err := json.Marshal(educationArray)
		if err != nil {
			log.Printf("❌ Failed to marshal education for member %d: %v", member.ID, err)
			continue
		}

		// Update database with proper JSON format
		if err := db.Table("members").Where("id = ?", member.ID).Update("education", string(educationJSON)).Error; err != nil {
			log.Printf("❌ Failed to update member %d: %v", member.ID, err)
			continue
		}

		fmt.Printf("   ✅ Successfully updated member ID %d\n", member.ID)
	}

	fmt.Println("\n🎉 Education migration completed!")
	fmt.Println("💡 All education fields are now in JSON array format")
}
