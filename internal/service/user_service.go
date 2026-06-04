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
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetAll(role string, isActive *bool, page, limit int) ([]*domain.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	users, total, err := s.repo.GetAll(role, isActive, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (s *userService) Create(username, email, password, fullName, role string) error {
	if username == "" || email == "" || password == "" || fullName == "" {
		return errors.New("missing required fields")
	}
	if role != "admin" && role != "teacher" {
		return errors.New("invalid role")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		Role:         role,
	}

	return s.repo.Create(user)
}

func (s *userService) Update(id, fullName, email, password string) error {
	if id == "" {
		return errors.New("missing user ID")
	}
	if fullName == "" || email == "" {
		return errors.New("missing required fields")
	}

	user := &domain.User{
		ID:       id,
		FullName: fullName,
		Email:    email,
	}

	if password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return err
		}
		user.PasswordHash = hashedPassword
	}

	return s.repo.Update(user)
}

func (s *userService) SoftDelete(id string) error {
	return s.repo.SoftDelete(id)
}
