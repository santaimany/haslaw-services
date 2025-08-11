package handlers

import (
	"haslaw-be-services/internal/service"
	"haslaw-be-services/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	ExpiresIn   int         `json:"expires_in"`
	User        interface{} `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Data request tidak valid", err.Error())
		return
	}

	user, token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.UnauthorizedResponse(c, "Login gagal")
		return
	}

	response := LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   900,
		User: map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Data request tidak valid", err.Error())
		return
	}

	newToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.UnauthorizedResponse(c, "Refresh token gagal")
		return
	}

	response := map[string]interface{}{
		"access_token": newToken,
		"token_type":   "Bearer",
		"expires_in":   900,
	}

	utils.SuccessResponse(c, http.StatusOK, "Token berhasil diperbarui", response)
}

func (h *AuthHandler) Logout(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Tidak diotorisasi")
		return
	}

	token, exists := c.Get("token")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak ditemukan")
		return
	}

	if err := h.authService.Logout(userID.(uint), token.(string)); err != nil {
		utils.InternalServerErrorResponse(c, "Logout gagal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Logout berhasil", nil)
}
