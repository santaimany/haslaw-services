package service

import (
	"haslaw-be-services/internal/models"
	"haslaw-be-services/internal/repository"
)

type MemberService interface {
	Create(memberData *CreateMemberRequest) (*models.Member, error)
	GetAll() ([]models.Member, error)
	GetByID(id uint) (*models.Member, error)
	Update(id uint, memberData *UpdateMemberRequest) (*models.Member, error)
	Delete(id uint) error
}

type CreateMemberRequest struct {
	FullName      string   `json:"full_name" binding:"required"`
	TitlePosition string   `json:"title_position" binding:"required"`
	Email         string   `json:"email" binding:"required,email"`
	PhoneNumber   string   `json:"phone_number"`
	LinkedIn      string   `json:"linked_in"`
	BusinessCard  string   `json:"business_card"`
	DisplayImage  string   `json:"display_image"`
	DetailImage   string   `json:"detail_image"`
	Biography     string   `json:"biography"`
	PracticeFocus []string `json:"practice_focus"`
	Education     []string `json:"education"`
	Language      []string `json:"language"`
}

type UpdateMemberRequest struct {
	FullName      string   `json:"full_name"`
	TitlePosition string   `json:"title_position"`
	Email         string   `json:"email"`
	PhoneNumber   string   `json:"phone_number"`
	LinkedIn      string   `json:"linked_in"`
	BusinessCard  string   `json:"business_card"`
	DisplayImage  string   `json:"display_image"`
	DetailImage   string   `json:"detail_image"`
	Biography     string   `json:"biography"`
	PracticeFocus []string `json:"practice_focus"`
	Education     []string `json:"education"`
	Language      []string `json:"language"`
}

type memberService struct {
	memberRepo repository.MemberRepository
}

func NewMemberService(memberRepo repository.MemberRepository) MemberService {
	return &memberService{
		memberRepo: memberRepo,
	}
}

func (s *memberService) Create(memberData *CreateMemberRequest) (*models.Member, error) {
	member := &models.Member{
		FullName:      memberData.FullName,
		TitlePosition: memberData.TitlePosition,
		Email:         memberData.Email,
		PhoneNumber:   memberData.PhoneNumber,
		LinkedIn:      memberData.LinkedIn,
		BusinessCard:  memberData.BusinessCard,
		DisplayImage:  memberData.DisplayImage,
		DetailImage:   memberData.DetailImage,
		Biography:     memberData.Biography,
		PracticeFocus: memberData.PracticeFocus,
		Education:     memberData.Education,
		Language:      memberData.Language,
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *memberService) GetAll() ([]models.Member, error) {
	return s.memberRepo.GetAll()
}

func (s *memberService) GetByID(id uint) (*models.Member, error) {
	return s.memberRepo.GetByID(id)
}

func (s *memberService) Update(id uint, memberData *UpdateMemberRequest) (*models.Member, error) {
	member, err := s.memberRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if memberData.FullName != "" {
		member.FullName = memberData.FullName
	}
	if memberData.TitlePosition != "" {
		member.TitlePosition = memberData.TitlePosition
	}
	if memberData.Email != "" {
		member.Email = memberData.Email
	}
	if memberData.PhoneNumber != "" {
		member.PhoneNumber = memberData.PhoneNumber
	}
	if memberData.LinkedIn != "" {
		member.LinkedIn = memberData.LinkedIn
	}
	if memberData.BusinessCard != "" {
		member.BusinessCard = memberData.BusinessCard
	}
	if memberData.DisplayImage != "" {
		member.DisplayImage = memberData.DisplayImage
	}
	if memberData.DetailImage != "" {
		member.DetailImage = memberData.DetailImage
	}
	if memberData.Biography != "" {
		member.Biography = memberData.Biography
	}
	if len(memberData.PracticeFocus) > 0 {
		member.PracticeFocus = memberData.PracticeFocus
	}
	if len(memberData.Education) > 0 {
		member.Education = memberData.Education
	}
	if len(memberData.Language) > 0 {
		member.Language = memberData.Language
	}

	if err := s.memberRepo.Update(member); err != nil {
		return nil, err
	}

	return member, nil
}

func (s *memberService) Delete(id uint) error {
	_, err := s.memberRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.memberRepo.Delete(id)
}
