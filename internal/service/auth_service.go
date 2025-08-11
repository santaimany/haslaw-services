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

type AuthService interface {
	Login(username, password string) (*models.User, string, error)
	CreateDefaultSuperAdmin() error
	CreateAdmin(request *models.CreateAdminRequest) (*models.User, error)
	UpdateProfile(userID uint, request *models.UpdateProfileRequest) (*models.User, error)
	ValidateToken(tokenString string) (*utils.Claims, error)
	RefreshToken(refreshToken string) (string, error)
	Logout(userID uint, token string) error
	IsTokenBlacklisted(token string) (bool, error)
}

type authService struct {
	userRepo      repository.UserRepository
	blacklistRepo repository.BlacklistRepository
}

func NewAuthService(userRepo repository.UserRepository, blacklistRepo repository.BlacklistRepository) AuthService {
	return &authService{
		userRepo:      userRepo,
		blacklistRepo: blacklistRepo,
	}
}

func (s *authService) Login(username, password string) (*models.User, string, error) {

	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	token, _, err := utils.GenerateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) CreateDefaultSuperAdmin() error {
	_, err := s.userRepo.GetByUsername("superadmin")
	if err == nil {

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

func (s *authService) ValidateToken(tokenString string) (*utils.Claims, error) {
	return utils.ValidateToken(tokenString)
}

func (s *authService) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	newToken, _, err := utils.GenerateTokens(claims.UserID, claims.Username, user.Role)
	return newToken, err
}

func (s *authService) CreateAdmin(request *models.CreateAdminRequest) (*models.User, error) {

	_, err := s.userRepo.GetByUsername(request.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	_, err = s.userRepo.GetByEmail(request.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

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

func (s *authService) UpdateProfile(userID uint, request *models.UpdateProfileRequest) (*models.User, error) {

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Username != request.Username {
		existingUser, err := s.userRepo.GetByUsername(request.Username)
		if err == nil && existingUser.ID != userID {
			return nil, errors.New("username already exists")
		}
	}

	if user.Email != request.Email {
		existingUser, err := s.userRepo.GetByEmail(request.Email)
		if err == nil && existingUser.ID != userID {
			return nil, errors.New("email already exists")
		}
	}

	user.Username = request.Username
	user.Email = request.Email

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

func (s *authService) Logout(userID uint, token string) error {

	claims, err := utils.ValidateToken(token)
	if err != nil {

		expiresAt := time.Now().Add(15 * time.Minute)
		return s.blacklistRepo.AddToBlacklist(token, userID, expiresAt)
	}

	expiresAt := time.Unix(claims.ExpiresAt.Unix(), 0)
	if err := s.blacklistRepo.AddToBlacklist(token, userID, expiresAt); err != nil {
		return err
	}

	return s.userRepo.UpdateRefreshToken(userID, "")
}

func (s *authService) IsTokenBlacklisted(token string) (bool, error) {
	return s.blacklistRepo.IsTokenBlacklisted(token)
}
