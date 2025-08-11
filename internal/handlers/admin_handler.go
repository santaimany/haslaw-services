package handlers

import (
	"haslaw-be-services/internal/models"
	"haslaw-be-services/internal/service"
	"haslaw-be-services/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	authService service.AuthService
}

func NewAdminHandler(authService service.AuthService) *AdminHandler {
	return &AdminHandler{
		authService: authService,
	}
}

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var request models.CreateAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	admin, err := h.authService.CreateAdmin(&request)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create admin", err.Error())
		return
	}

	response := models.UserResponse{
		ID:       admin.ID,
		Username: admin.Username,
		Email:    admin.Email,
		Role:     admin.Role,
	}

	utils.SuccessResponse(c, http.StatusCreated, "Admin created successfully", response)
}

func (h *AdminHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User ID not found in token")
		return
	}

	var request models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, err := h.authService.UpdateProfile(userID.(uint), &request)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to update profile", err.Error())
		return
	}

	response := models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", response)
}

func (h *AdminHandler) GetProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	if userIDStr == "" {

		userID, exists := c.Get("user_id")
		if !exists {
			utils.UnauthorizedResponse(c, "User ID not found in token")
			return
		}
		userIDStr = strconv.Itoa(int(userID.(uint)))
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID", err.Error())
		return
	}

	currentUserID, _ := c.Get("user_id")
	currentUserRole, _ := c.Get("role")

	if currentUserID.(uint) != uint(userID) && currentUserRole != string(models.SuperAdmin) {
		utils.ForbiddenResponse(c, "You can only access your own profile")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", map[string]interface{}{
		"message": "Profile endpoint - implementation can be expanded",
		"user_id": userID,
	})
}
