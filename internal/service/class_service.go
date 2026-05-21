package service

import (
	"alhikmah-attendance-api/internal/domain"
)

type classService struct {
	repo domain.ClassRepository
}

func NewClassService(repo domain.ClassRepository) domain.ClassService {
	return &classService{repo: repo}
}

func (s *classService) GetAll(teacherID string, academicYear string, page, limit int) ([]*domain.Class, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.repo.GetAll(teacherID, academicYear, offset, limit)
}

func (s *classService) GetByID(id string) (*domain.Class, error) {
	return s.repo.GetByID(id)
}
