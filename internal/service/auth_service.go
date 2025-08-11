package service

import (
	"errors"
	"haslaw-be-services/internal/models"
	"haslaw-be-services/internal/repository"
	"haslaw-be-services/internal/utils"
	"log"
	"time"

	"gorm.io/gorm"
)

// AuthService interface defines authentication business logic
type AuthService interface {
	Login(username, password string) (*models.User, string, string, error) // Returns user, accessToken, refreshToken, error
	CreateDefaultSuperAdmin() error
	CreateAdmin(request *models.CreateAdminRequest) (*models.User, error)
	UpdateProfile(userID uint, request *models.UpdateProfileRequest) (*models.User, error)
	ValidateToken(tokenString string) (*utils.Claims, error)
	RefreshToken(refreshToken string) (string, string, error) // Returns newAccessToken, newRefreshToken, error
	Logout(userID uint, token string) error
	IsTokenBlacklisted(token string) (bool, error)
}

// authService implements AuthService interface
type authService struct {
	userRepo      repository.UserRepository
	blacklistRepo repository.BlacklistRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, blacklistRepo repository.BlacklistRepository) AuthService {
	return &authService{
		userRepo:      userRepo,
		blacklistRepo: blacklistRepo,
	}
}

// Login authenticates user and returns JWT tokens
func (s *authService) Login(username, password string) (*models.User, string, string, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("invalid credentials")
		}
		return nil, "", "", err
	}

	// Check password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", "", errors.New("invalid credentials")
	}

	// Generate both access and refresh tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

// CreateDefaultSuperAdmin creates default super admin user if not exists
func (s *authService) CreateDefaultSuperAdmin() error {
	_, err := s.userRepo.GetByUsername("superadmin")
	if err == nil {
		// Super admin already exists
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create default super admin user
	hashedPassword, err := utils.HashPassword("superadmin123")
	if err != nil {
		return err
	}

	superAdmin := &models.User{
		Username: "superadmin",
		Email:    "superadmin@haslaw.com",
		Password: hashedPassword,
		Role:     models.SuperAdmin,
	}

	if err := s.userRepo.Create(superAdmin); err != nil {
		return err
	}

	log.Println("Default super admin user created:")
	log.Println("Username: superadmin")
	log.Println("Password: superadmin123")
	log.Println("Please change the password after first login!")

	return nil
}

// ValidateToken validates JWT token and returns claims
func (s *authService) ValidateToken(tokenString string) (*utils.Claims, error) {
	return utils.ValidateToken(tokenString)
}

// RefreshToken generates new access token and refresh token from old refresh token
func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Check if refresh token is blacklisted
	isBlacklisted, err := s.IsTokenBlacklisted(refreshToken)
	if err != nil {
		return "", "", errors.New("failed to check token status")
	}
	if isBlacklisted {
		return "", "", errors.New("refresh token is blacklisted")
	}

	// Get user to get current role
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	// Blacklist the old refresh token
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 hari
	if err := s.blacklistRepo.AddToBlacklist(refreshToken, claims.UserID, expiresAt); err != nil {
		log.Printf("Warning: Failed to blacklist old refresh token: %v", err)
	}

	// Generate new tokens (access token: 15 menit, refresh token: 7 hari)
	newAccessToken, newRefreshToken, err := utils.GenerateTokens(claims.UserID, claims.Username, user.Role)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// CreateAdmin creates a new admin user (only super admin can do this)
func (s *authService) CreateAdmin(request *models.CreateAdminRequest) (*models.User, error) {
	// Check if username already exists
	_, err := s.userRepo.GetByUsername(request.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	_, err = s.userRepo.GetByEmail(request.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	// Create new admin user
	admin := &models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
		Role:     models.Admin,
	}

	if err := s.userRepo.Create(admin); err != nil {
		return nil, err
	}

	return admin, nil
}

// UpdateProfile updates user profile (admin can only update their own profile)
func (s *authService) UpdateProfile(userID uint, request *models.UpdateProfileRequest) (*models.User, error) {
	// Get current user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if username is taken by another user
	if user.Username != request.Username {
		existingUser, err := s.userRepo.GetByUsername(request.Username)
		if err == nil && existingUser.ID != userID {
			return nil, errors.New("username already exists")
		}
	}

	// Check if email is taken by another user
	if user.Email != request.Email {
		existingUser, err := s.userRepo.GetByEmail(request.Email)
		if err == nil && existingUser.ID != userID {
			return nil, errors.New("email already exists")
		}
	}

	// Update user fields
	user.Username = request.Username
	user.Email = request.Email

	// Update password if provided
	if request.Password != "" {
		hashedPassword, err := utils.HashPassword(request.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Logout invalidates user session by blacklisting the token
func (s *authService) Logout(userID uint, token string) error {
	// Parse token to get expiry time
	claims, err := utils.ValidateToken(token)
	if err != nil {
		// Even if token is invalid, we should still try to blacklist it
		// Use a default expiry time if we can't parse the token
		expiresAt := time.Now().Add(15 * time.Minute)
		return s.blacklistRepo.AddToBlacklist(token, userID, expiresAt)
	}

	// Add token to blacklist with its original expiry time
	expiresAt := time.Unix(claims.ExpiresAt.Unix(), 0)
	if err := s.blacklistRepo.AddToBlacklist(token, userID, expiresAt); err != nil {
		return err
	}

	// Also clear refresh token from user record
	return s.userRepo.UpdateRefreshToken(userID, "")
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *authService) IsTokenBlacklisted(token string) (bool, error) {
	return s.blacklistRepo.IsTokenBlacklisted(token)
}
