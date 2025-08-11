package handlers

import (
	"haslaw-be-services/internal/service"
	"haslaw-be-services/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler untuk menghandle permintaan autentikasi
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler membuat auth handler baru
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest untuk data login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	ExpiresIn   int         `json:"expires_in"`
	User        interface{} `json:"user"`
	Message     string      `json:"message"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	// Initialize trace logger
	logger := utils.NewTraceLogger(c)
	logger.Info("Login attempt started")

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err, "Invalid login request payload")
		utils.BadRequestResponse(c, "Data request tidak valid", err.Error())
		return
	}

	logger.Infof("Login attempt for username: %s", req.Username)

	user, accessToken, refreshToken, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		logger.WithError(err, "Login authentication failed")
		utils.UnauthorizedResponse(c, "Login gagal")
		return
	}

	logger.LogUserAction(user.Username, "LOGIN", "auth")

	c.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*60*60,
		"/",
		"",
		false,
		true,
	)

	// Buat response login (tanpa refresh token di body, karena sudah di cookie)
	response := LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   900, // 15 menit untuk access token
		User: map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
		Message: "Refresh token tersimpan di cookie (7 hari)",
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", response)
}

// RefreshToken untuk memperbarui token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Ambil refresh token dari cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.UnauthorizedResponse(c, "Refresh token tidak ditemukan di cookie")
		return
	}

	// Proses refresh token
	newAccessToken, newRefreshToken, err := h.authService.RefreshToken(refreshToken)
	if err != nil {
		// Hapus cookie jika refresh token tidak valid
		c.SetCookie("refresh_token", "", -1, "/", "", false, true)
		utils.UnauthorizedResponse(c, "Refresh token tidak valid atau expired")
		return
	}

	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		7*24*60*60,
		"/",
		"",
		false,
		true,
	)

	response := map[string]interface{}{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
		"expires_in":   900,
		"message":      "Token berhasil diperbarui, refresh token baru tersimpan di cookie",
	}

	utils.SuccessResponse(c, http.StatusOK, "Token berhasil diperbarui", response)
}

func (h *AuthHandler) Logout(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "Tidak diotorisasi")
		return
	}

	accessToken, exists := c.Get("token")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak ditemukan")
		return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {

		log.Printf("Warning: No refresh token found in cookie during logout: %v", err)
	}

	if err := h.authService.Logout(userID.(uint), accessToken.(string)); err != nil {
		utils.InternalServerErrorResponse(c, "Logout gagal", err.Error())
		return
	}

	if refreshToken != "" {
		log.Printf("Refresh token removed from cookie during logout")
	}

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	utils.SuccessResponse(c, http.StatusOK, "Logout berhasil, semua token telah dihapus", nil)
}
