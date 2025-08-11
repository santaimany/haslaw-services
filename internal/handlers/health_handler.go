package handlers

import (
	"haslaw-be-services/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	healthData := map[string]interface{}{
		"status":  "healthy",
		"service": "haslaw-services",
		"version": "1.0.0",
	}

	utils.SuccessResponse(c, http.StatusOK, "Service is healthy", healthData)
}
