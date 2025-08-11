package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"unique;not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	Password     string    `json:"-" gorm:"not null"`                                     // Password tidak ditampilkan di JSON
	Role         UserRole  `json:"role" gorm:"type:varchar(20);not null;default:'admin'"` // Role user
	RefreshToken string    `json:"-"`                                                     // Refresh token untuk login
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type News struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	NewsTitle string         `json:"news_title" gorm:"not null"`                                // Judul berita
	Slug      string         `json:"slug" gorm:"unique;not null;index"`                         // URL slug untuk berita
	Category  string         `json:"category" gorm:"not null"`                                  // Kategori berita
	Status    NewsStatus     `json:"status" gorm:"type:varchar(20);not null;default:'Drafted'"` // Status publish
	Content   string         `json:"content" gorm:"type:text"`                                  // Isi berita
	Image     string         `json:"image"`                                                     // Gambar berita
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete
}

type Member struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	FullName      string         `json:"full_name" gorm:"not null"`             // Nama lengkap
	TitlePosition string         `json:"title_position" gorm:"not null"`        // Jabatan
	Email         string         `json:"email" gorm:"unique;not null"`          // Email member
	PhoneNumber   string         `json:"phone_number"`                          // Nomor telepon
	LinkedIn      string         `json:"linkedin"`                              // Profile LinkedIn
	BusinessCard  string         `json:"business_card"`                         // File kartu nama
	DisplayImage  string         `json:"display_image"`                         // Foto profil
	DetailImage   string         `json:"detail_image"`                          // Foto detail
	Biography     string         `json:"biography" gorm:"type:text"`            // Biografi
	PracticeFocus []string       `json:"practice_focus" gorm:"serializer:json"` // Bidang keahlian
	Education     string         `json:"education"`                             // Pendidikan
	Language      []string       `json:"language" gorm:"serializer:json"`       // Bahasa yang dikuasai
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete
}

type BlacklistedToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"not null;type:text"`                   // JWT token yang di-blacklist (full token)
	TokenHash string    `json:"token_hash" gorm:"not null;index;type:varchar(64)"` // SHA256 hash of token for indexing
	UserID    uint      `json:"user_id" gorm:"not null"`                           // ID user yang logout
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index"`                  // Kapan token expire
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     UserRole `json:"role"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

type CreateAdminRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
}

type NewsRequest struct {
	NewsTitle string `json:"news_title" binding:"required"`
	Slug      string `json:"slug"`
	Category  string `json:"category" binding:"required"`
	Status    string `json:"status"`
	Content   string `json:"content"`
}

type MemberRequest struct {
	FullName      string   `json:"full_name" binding:"required"`
	TitlePosition string   `json:"title_position" binding:"required"`
	Email         string   `json:"email" binding:"required,email"`
	PhoneNumber   string   `json:"phone_number"`
	LinkedIn      string   `json:"linkedin"`
	BusinessCard  string   `json:"business_card"`
	Biography     string   `json:"biography"`
	PracticeFocus []string `json:"practice_focus"`
	Education     string   `json:"education"`
	Language      []string `json:"language"`
}
