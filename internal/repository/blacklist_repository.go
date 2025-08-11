package repository

import (
	"haslaw-be-services/internal/models"
	"time"

	"gorm.io/gorm"
)

type BlacklistRepository interface {
	AddToBlacklist(token string, userID uint, expiresAt time.Time) error
	IsTokenBlacklisted(token string) (bool, error)
	CleanupExpiredTokens() error
}

type blacklistRepository struct {
	db *gorm.DB
}

func NewBlacklistRepository(db *gorm.DB) BlacklistRepository {
	return &blacklistRepository{db: db}
}

func (r *blacklistRepository) AddToBlacklist(token string, userID uint, expiresAt time.Time) error {
	blacklistedToken := &models.BlacklistedToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	return r.db.Create(blacklistedToken).Error
}

func (r *blacklistRepository) IsTokenBlacklisted(token string) (bool, error) {
	var count int64
	err := r.db.Model(&models.BlacklistedToken{}).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *blacklistRepository) CleanupExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.BlacklistedToken{}).Error
}
