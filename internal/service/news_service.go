package service

import (
	"errors"
	"haslaw-be-services/internal/models"
	"haslaw-be-services/internal/repository"
	"haslaw-be-services/internal/utils"

	"gorm.io/gorm"
)

type NewsService interface {
	Create(newsData *CreateNewsRequest) (*models.News, error)
	GetAll(page, limit int, orderBy, category string) ([]models.News, *utils.PaginationMeta, error)
	GetPublished(page, limit int, orderBy, category string) ([]models.News, *utils.PaginationMeta, error)
	GetDrafts(page, limit int, orderBy string) ([]models.News, *utils.PaginationMeta, error)
	GetByID(id uint) (*models.News, error)
	GetBySlug(slug string) (*models.News, error)
	Update(id uint, newsData *UpdateNewsRequest) (*models.News, error)
	Delete(id uint) error
	Publish(id uint) (*models.News, error)
}

type CreateNewsRequest struct {
	NewsTitle string            `json:"news_title" binding:"required"`
	Category  string            `json:"category" binding:"required"`
	Status    models.NewsStatus `json:"status" binding:"required"`
	Content   string            `json:"content"`
	Image     string            `json:"image"`
}

type UpdateNewsRequest struct {
	NewsTitle string            `json:"news_title"`
	Category  string            `json:"category"`
	Status    models.NewsStatus `json:"status"`
	Content   string            `json:"content"`
	Image     string            `json:"image"`
}

type newsService struct {
	newsRepo repository.NewsRepository
}

func NewNewsService(newsRepo repository.NewsRepository) NewsService {
	return &newsService{
		newsRepo: newsRepo,
	}
}

func (s *newsService) Create(newsData *CreateNewsRequest) (*models.News, error) {

	if !newsData.Status.IsValid() {
		return nil, errors.New("invalid news status")
	}

	slug := utils.GenerateSlugWithRandomID(newsData.NewsTitle)

	news := &models.News{
		NewsTitle: newsData.NewsTitle,
		Slug:      slug,
		Category:  newsData.Category,
		Status:    newsData.Status,
		Content:   newsData.Content,
		Image:     newsData.Image,
	}

	if err := s.newsRepo.Create(news); err != nil {
		return nil, err
	}

	return news, nil
}

func (s *newsService) GetAll(page, limit int, orderBy, category string) ([]models.News, *utils.PaginationMeta, error) {
	offset := (page - 1) * limit
	orderClause := s.buildOrderClause(orderBy)

	news, total, err := s.newsRepo.GetAll(limit, offset, orderClause)
	if err != nil {
		return nil, nil, err
	}

	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	}

	return news, meta, nil
}

func (s *newsService) GetPublished(page, limit int, orderBy, category string) ([]models.News, *utils.PaginationMeta, error) {
	offset := (page - 1) * limit
	orderClause := s.buildOrderClause(orderBy)

	news, total, err := s.newsRepo.GetPublished(limit, offset, orderClause, category)
	if err != nil {
		return nil, nil, err
	}

	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	}

	return news, meta, nil
}

func (s *newsService) GetDrafts(page, limit int, orderBy string) ([]models.News, *utils.PaginationMeta, error) {
	offset := (page - 1) * limit
	orderClause := s.buildOrderClause(orderBy)

	news, total, err := s.newsRepo.GetDrafts(limit, offset, orderClause)
	if err != nil {
		return nil, nil, err
	}

	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: (total + int64(limit) - 1) / int64(limit),
	}

	return news, meta, nil
}

func (s *newsService) GetByID(id uint) (*models.News, error) {
	return s.newsRepo.GetByID(id)
}

func (s *newsService) GetBySlug(slug string) (*models.News, error) {
	return s.newsRepo.GetBySlug(slug)
}

func (s *newsService) Update(id uint, newsData *UpdateNewsRequest) (*models.News, error) {
	news, err := s.newsRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("news not found")
		}
		return nil, err
	}

	if newsData.NewsTitle != "" {
		news.NewsTitle = newsData.NewsTitle

		news.Slug = utils.GenerateSlugWithRandomID(newsData.NewsTitle)
	}
	if newsData.Category != "" {
		news.Category = newsData.Category
	}
	if newsData.Status != "" {
		if !newsData.Status.IsValid() {
			return nil, errors.New("invalid news status")
		}
		news.Status = newsData.Status
	}
	if newsData.Content != "" {
		news.Content = newsData.Content
	}
	if newsData.Image != "" {
		news.Image = newsData.Image
	}

	if err := s.newsRepo.Update(news); err != nil {
		return nil, err
	}

	return news, nil
}

func (s *newsService) Delete(id uint) error {
	_, err := s.newsRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("news not found")
		}
		return err
	}

	return s.newsRepo.Delete(id)
}

func (s *newsService) Publish(id uint) (*models.News, error) {
	news, err := s.newsRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("news not found")
		}
		return nil, err
	}

	if news.Status != models.Drafted {
		return nil, errors.New("only draft news can be published")
	}

	if err := s.newsRepo.Publish(id); err != nil {
		return nil, err
	}

	return s.newsRepo.GetByID(id)
}

func (s *newsService) buildOrderClause(orderBy string) string {
	switch orderBy {
	case "id_asc":
		return "id ASC"
	case "id_desc":
		return "id DESC"
	case "title_asc":
		return "news_title ASC"
	case "title_desc":
		return "news_title DESC"
	case "created_at_asc":
		return "created_at ASC"
	case "created_at_desc":
		return "created_at DESC"
	case "updated_at_asc":
		return "updated_at ASC"
	case "updated_at_desc":
		return "updated_at DESC"
	default:
		return "created_at DESC"
	}
}
