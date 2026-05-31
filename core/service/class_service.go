package service

import (
	"errors"

	"alhikmah-attendance-api/core/domain"
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

func (s *classService) Create(class *domain.Class) error {
	if class.ClassName == "" {
		return errors.New("class_name is required")
	}
	if class.TeacherID == "" {
		return errors.New("teacher_id is required")
	}
	if class.AcademicYear == "" || len(class.AcademicYear) != 9 || class.AcademicYear[4] != '/' {
		return errors.New("academic_year is required and must be in format YYYY/YYYY")
	}
	if class.Capacity <= 0 {
		class.Capacity = 30
	}
	return s.repo.Create(class)
}

func (s *classService) Update(class *domain.Class) error {
	if class.ID == "" {
		return errors.New("class_id is required")
	}
	if class.ClassName == "" {
		return errors.New("class_name is required")
	}
	if class.TeacherID == "" {
		return errors.New("teacher_id is required")
	}
	if class.Capacity <= 0 {
		class.Capacity = 30
	}
	return s.repo.Update(class)
}

func (s *classService) SoftDelete(id string) error {
	class, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if class.StudentCount > 0 {
		return errors.New("cannot delete class with active students")
	}
	return s.repo.SoftDelete(id)
}
