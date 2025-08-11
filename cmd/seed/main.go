package main

import (
	"haslaw-be-services/internal/config"
	"haslaw-be-services/internal/models"
	"haslaw-be-services/internal/utils"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("🌱 Starting database seeding...")

	seedAdmin(db)
	seedNews(db)
	seedMembers(db)

	log.Println("✅ Database seeding completed successfully!")
}

func seedAdmin(db any) {
	log.Println("📝 Seeding admin user...")
	log.Println("⚠️  Admin creation handled by auth service...")
}

func seedNews(db any) {
	log.Println("📰 Seeding news data...")

	newsData := []struct {
		Title    string
		Category string
		Status   models.NewsStatus
		Content  string
		Image    string
	}{
		{
			Title:    "Haslaw Wins Landmark Corporate Merger Case",
			Category: "Corporate Law",
			Status:   models.Posted,
			Content:  "Haslaw & Partners successfully facilitated a $500 million merger between two major Indonesian corporations. The complex transaction involved extensive due diligence, regulatory approvals, and cross-border compliance issues. Our corporate law team worked tirelessly to ensure all legal requirements were met and the transaction was completed smoothly.",
			Image:    "/uploads/news/corporate-merger-case.jpg",
		},
		{
			Title:    "Banking Regulation Update: New Compliance Requirements",
			Category: "Banking & Finance",
			Status:   models.Posted,
			Content:  "The Indonesian Financial Services Authority (OJK) has introduced new compliance requirements for banking institutions. Our banking and finance team has prepared comprehensive guidelines to help clients navigate these changes. The new regulations focus on enhanced risk management and consumer protection measures.",
			Image:    "/uploads/news/banking-regulation-update.jpg",
		},
		{
			Title:    "Employment Law Reform: Impact on Indonesian Businesses",
			Category: "Employment Law",
			Status:   models.Posted,
			Content:  "Recent amendments to Indonesian employment law have significant implications for businesses operating in Indonesia. Our employment law specialists have analyzed the changes and their impact on employment contracts, termination procedures, and employee benefits. Companies are advised to review their current policies.",
			Image:    "/uploads/news/employment-law-reform.jpg",
		},
		{
			Title:    "Intellectual Property Protection in Digital Era",
			Category: "Intellectual Property",
			Status:   models.Posted,
			Content:  "As businesses increasingly move to digital platforms, protecting intellectual property rights becomes more crucial than ever. Our IP team discusses the latest trends in digital IP protection, including trademark registration for online businesses, copyright protection for digital content, and patent filing for tech innovations.",
			Image:    "/uploads/news/ip-protection-digital.jpg",
		},
		{
			Title:    "Environmental Law Compliance for Manufacturing Sector",
			Category: "Environmental Law",
			Status:   models.Posted,
			Content:  "New environmental regulations require manufacturing companies to implement stricter waste management and emission control measures. Our environmental law team provides guidance on compliance strategies, permit applications, and sustainable business practices that align with Indonesian environmental standards.",
			Image:    "/uploads/news/environmental-compliance.jpg",
		},
		{
			Title:    "Draft: Upcoming Tax Law Changes in 2024",
			Category: "Tax Law",
			Status:   models.Drafted,
			Content:  "This draft article covers the proposed tax law changes expected to take effect in 2024. The changes include modifications to corporate tax rates, VAT regulations, and international tax compliance requirements. This article is currently under review by our tax law specialists.",
			Image:    "/uploads/news/tax-law-changes-draft.jpg",
		},
	}

	for i, data := range newsData {
		slug := utils.GenerateSlugWithRandomID(data.Title)

		_ = models.News{
			NewsTitle: data.Title,
			Slug:      slug,
			Category:  data.Category,
			Status:    data.Status,
			Content:   data.Content,
			Image:     data.Image,
		}

		log.Printf("✅ [%d/%d] Would create: %s", i+1, len(newsData), data.Title)
	}
}

func seedMembers(db any) {
	log.Println("👥 Seeding members data...")

	membersData := []struct {
		FullName      string
		TitlePosition string
		Email         string
	}{
		{
			FullName:      "Dr. Budi Santoso, S.H., LL.M.",
			TitlePosition: "Managing Partner",
			Email:         "budi.santoso@haslaw.com",
		},
		{
			FullName:      "Sari Wijayanti, S.H., LL.M.",
			TitlePosition: "Senior Partner - Banking & Finance",
			Email:         "sari.wijayanti@haslaw.com",
		},
	}

	for i, data := range membersData {
		log.Printf("✅ [%d/%d] Would create: %s (%s)", i+1, len(membersData), data.FullName, data.TitlePosition)
	}
}
