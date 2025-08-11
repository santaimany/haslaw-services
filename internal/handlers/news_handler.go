package handlers

import (
	"fmt"
	"haslaw-be-services/internal/models"
	"haslaw-be-services/internal/service"
	"haslaw-be-services/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// NewsHandler handles news requests
type NewsHandler struct {
	newsService service.NewsService
}

// NewNewsHandler creates a new news handler
func NewNewsHandler(newsService service.NewsService) *NewsHandler {
	return &NewsHandler{
		newsService: newsService,
	}
}

// GetAllPublicNews gets all published news for public
func (h *NewsHandler) GetAllPublicNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orderBy := c.DefaultQuery("order_by", "created_at_desc")
	category := c.Query("category")

	news, meta, err := h.newsService.GetPublished(page, limit, orderBy, category)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch news", err.Error())
		return
	}

	utils.SuccessWithPagination(c, "News retrieved successfully", news, *meta)
}

// GetAll gets all news for admin
func (h *NewsHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orderBy := c.DefaultQuery("order_by", "created_at_desc")
	category := c.Query("category")

	news, meta, err := h.newsService.GetAll(page, limit, orderBy, category)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch news", err.Error())
		return
	}

	utils.SuccessWithPagination(c, "News retrieved successfully", news, *meta)
}

// GetByID gets news by ID
func (h *NewsHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid news ID", err.Error())
		return
	}

	news, err := h.newsService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "News not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "News retrieved successfully", news)
}

// GetBySlug gets news by slug
func (h *NewsHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	news, err := h.newsService.GetBySlug(slug)
	if err != nil {
		utils.NotFoundResponse(c, "News not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "News retrieved successfully", news)
}

// Create creates new news
func (h *NewsHandler) Create(c *gin.Context) {
	var req service.CreateNewsRequest
	
	// Check content type - support both JSON and form-data
	contentType := c.GetHeader("Content-Type")
	
	if strings.Contains(contentType, "application/json") {
		// Handle JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequestResponse(c, "Invalid request body", err.Error())
			return
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		// Handle form-data request
		req.NewsTitle = c.PostForm("news_title")
		req.Category = c.PostForm("category")
		req.Status = models.NewsStatus(c.PostForm("status"))
		req.Content = c.PostForm("content")
		
		// Handle file upload
		file, err := c.FormFile("image")
		if err != nil {
			utils.BadRequestResponse(c, "Image file is required", err.Error())
			return
		}
		
		// Validate required fields
		if req.NewsTitle == "" || req.Category == "" || req.Status == "" || req.Content == "" {
			utils.BadRequestResponse(c, "All fields are required", "news_title, category, status, and content are required")
			return
		}
		
		// Save uploaded file
		uploadPath := fmt.Sprintf("uploads/news/%d_%s", time.Now().Unix(), file.Filename)
		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			utils.InternalServerErrorResponse(c, "Failed to save image", err.Error())
			return
		}
		
		req.Image = uploadPath
	} else {
		utils.BadRequestResponse(c, "Unsupported content type", "Use application/json or multipart/form-data")
		return
	}

	news, err := h.newsService.Create(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create news", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "News created successfully", news)
}

// Update updates news
func (h *NewsHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid news ID", err.Error())
		return
	}

	var req service.UpdateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err.Error())
		return
	}

	news, err := h.newsService.Update(uint(id), &req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to update news", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "News updated successfully", news)
}

// Delete deletes news
func (h *NewsHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid news ID", err.Error())
		return
	}

	if err := h.newsService.Delete(uint(id)); err != nil {
		utils.BadRequestResponse(c, "Failed to delete news", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "News deleted successfully", nil)
}

// GetDrafts gets draft news
func (h *NewsHandler) GetDrafts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orderBy := c.DefaultQuery("order_by", "created_at_desc")

	news, meta, err := h.newsService.GetDrafts(page, limit, orderBy)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch draft news", err.Error())
		return
	}

	utils.SuccessWithPagination(c, "Draft news retrieved successfully", news, *meta)
}

// GetDraftByID gets draft news by ID
func (h *NewsHandler) GetDraftByID(c *gin.Context) {
	h.GetByID(c) // Same logic as GetByID
}

// PublishDraft publishes a draft news
func (h *NewsHandler) PublishDraft(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid news ID", err.Error())
		return
	}

	news, err := h.newsService.Publish(uint(id))
	if err != nil {
		utils.BadRequestResponse(c, "Failed to publish news", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "News published successfully", news)
}
