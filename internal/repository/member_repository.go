package repository

import (
	"haslaw-be-services/internal/models"

	"gorm.io/gorm"
)

type MemberRepository interface {
	Create(member *models.Member) error
	GetAll() ([]models.Member, error)
	GetByID(id uint) (*models.Member, error)
	GetByEmail(email string) (*models.Member, error)
	Update(member *models.Member) error
	Delete(id uint) error
}

type memberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{db: db}
}

func (r *memberRepository) Create(member *models.Member) error {
	return r.db.Create(member).Error
}

func (r *memberRepository) GetAll() ([]models.Member, error) {
	var members []models.Member
	err := r.db.Find(&members).Error
	return members, err
}

func (r *memberRepository) GetByID(id uint) (*models.Member, error) {
	var member models.Member
	err := r.db.First(&member, id).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) GetByEmail(email string) (*models.Member, error) {
	var member models.Member
	err := r.db.Where("email = ?", email).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) Update(member *models.Member) error {
	return r.db.Save(member).Error
}

func (r *memberRepository) Delete(id uint) error {
	return r.db.Delete(&models.Member{}, id).Error
}
