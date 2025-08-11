package repository

import (
	"haslaw-be-services/internal/models"

	"gorm.io/gorm"
)

type NewsRepository interface {
	Create(news *models.News) error
	GetAll(limit, offset int, orderBy string) ([]models.News, int64, error)
	GetPublished(limit, offset int, orderBy string, category string) ([]models.News, int64, error)
	GetDrafts(limit, offset int, orderBy string) ([]models.News, int64, error)
	GetByID(id uint) (*models.News, error)
	GetBySlug(slug string) (*models.News, error)
	Update(news *models.News) error
	Delete(id uint) error
	Publish(id uint) error
	GetByCategory(category string, limit, offset int) ([]models.News, int64, error)
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) Create(news *models.News) error {
	return r.db.Create(news).Error
}

func (r *newsRepository) GetAll(limit, offset int, orderBy string) ([]models.News, int64, error) {
	var news []models.News
	var total int64

	// Use single query with count estimation for better performance
	query := r.db.Model(&models.News{})

	// Perform count and select in parallel-like manner
	countChan := make(chan error, 1)
	go func() {
		countChan <- query.Count(&total).Error
	}()

	// Execute main query with optimizations - include content field
	err := r.db.Select("id, news_title, slug, category, status, content, image, created_at, updated_at").
		Offset(offset).
		Limit(limit).
		Order(orderBy).
		Find(&news).Error

	// Wait for count query
	if countErr := <-countChan; countErr != nil {
		return nil, 0, countErr
	}

	return news, total, err
}

func (r *newsRepository) GetPublished(limit, offset int, orderBy string, category string) ([]models.News, int64, error) {
	var news []models.News
	var total int64

	baseQuery := r.db.Model(&models.News{}).Where("status = ?", models.Posted)

	if category != "" {
		baseQuery = baseQuery.Where("category = ?", category)
	}

	// Parallel count and select
	countChan := make(chan error, 1)
	go func() {
		countChan <- baseQuery.Count(&total).Error
	}()

	// Optimized select query with limited fields for list view - include content
	selectQuery := r.db.Select("id, news_title, slug, category, status, content, image, created_at, updated_at").
		Where("status = ?", models.Posted)

	if category != "" {
		selectQuery = selectQuery.Where("category = ?", category)
	}

	err := selectQuery.Offset(offset).
		Limit(limit).
		Order(orderBy).
		Find(&news).Error

	// Wait for count
	if countErr := <-countChan; countErr != nil {
		return nil, 0, countErr
	}

	return news, total, err
}

func (r *newsRepository) GetDrafts(limit, offset int, orderBy string) ([]models.News, int64, error) {
	var news []models.News
	var total int64

	query := r.db.Model(&models.News{}).Where("status = ?", models.Drafted)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order(orderBy).Find(&news).Error; err != nil {
		return nil, 0, err
	}

	return news, total, nil
}

func (r *newsRepository) GetByID(id uint) (*models.News, error) {
	var news models.News
	err := r.db.First(&news, id).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *newsRepository) GetBySlug(slug string) (*models.News, error) {
	var news models.News
	err := r.db.Where("slug = ?", slug).First(&news).Error
	if err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *newsRepository) Update(news *models.News) error {
	return r.db.Save(news).Error
}

func (r *newsRepository) Delete(id uint) error {
	return r.db.Delete(&models.News{}, id).Error
}

func (r *newsRepository) Publish(id uint) error {
	return r.db.Model(&models.News{}).Where("id = ?", id).Update("status", models.Posted).Error
}

func (r *newsRepository) GetByCategory(category string, limit, offset int) ([]models.News, int64, error) {
	var news []models.News
	var total int64

	query := r.db.Model(&models.News{}).Where("category = ? AND status = ?", category, models.Posted)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&news).Error; err != nil {
		return nil, 0, err
	}

	return news, total, nil
}
