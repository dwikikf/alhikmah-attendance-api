package service

import (
	"errors"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/pkg/utils"
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) GetByID(id string) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) GetAll(role string, isActive *bool, page, limit int) ([]*domain.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	return s.repo.GetAll(role, isActive, offset, limit)
}

func (s *userService) Create(user *domain.User) error {
	// Validate inputs
	if user.Username == "" || user.Email == "" || user.PasswordHash == "" || user.FullName == "" {
		return errors.New("missing required fields")
	}

	if user.Role != "admin" && user.Role != "teacher" {
		return errors.New("invalid role")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.PasswordHash) // Assuming raw password is in PasswordHash temporarily
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	return s.repo.Create(user)
}

func (s *userService) Update(user *domain.User) error {
	if user.ID == "" {
		return errors.New("missing user ID")
	}
	if user.FullName == "" || user.Email == "" {
		return errors.New("missing required fields")
	}

	return s.repo.Update(user)
}

func (s *userService) SoftDelete(id string) error {
	return s.repo.SoftDelete(id)
}
