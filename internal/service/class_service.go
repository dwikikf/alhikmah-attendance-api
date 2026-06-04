package service

import (
	"errors"

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
	classes, total, err := s.repo.GetAll(teacherID, academicYear, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

func (s *classService) GetByID(id string) (*domain.Class, error) {
	class, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return class, nil
}

func (s *classService) Create(roomName string, grade int, section *int, teacherID, academicYear string, capacity int, description string) error {
	if roomName == "" {
		return errors.New("room_name is required")
	}
	if grade < 1 || grade > 12 {
		return errors.New("grade must be between 1 and 12")
	}
	if section != nil && *section <= 0 {
		return errors.New("section must be greater than 0")
	}
	if teacherID == "" {
		return errors.New("teacher_id is required")
	}
	if academicYear == "" || len(academicYear) != 9 || academicYear[4] != '/' {
		return errors.New("academic_year is required and must be in format YYYY/YYYY")
	}
	if capacity <= 0 {
		capacity = 30
	}

	class := &domain.Class{
		RoomName:     roomName,
		Grade:        grade,
		Section:      section,
		TeacherID:    teacherID,
		AcademicYear: academicYear,
		Capacity:     capacity,
		Description:  description,
	}

	return s.repo.Create(class)
}

func (s *classService) Update(id, roomName string, grade int, section *int, teacherID string, capacity int, description string) error {
	if id == "" {
		return errors.New("class_id is required")
	}
	if roomName == "" {
		return errors.New("room_name is required")
	}
	if grade < 1 || grade > 12 {
		return errors.New("grade must be between 1 and 12")
	}
	if section != nil && *section <= 0 {
		return errors.New("section must be greater than 0")
	}
	if teacherID == "" {
		return errors.New("teacher_id is required")
	}
	if capacity <= 0 {
		capacity = 30
	}

	class := &domain.Class{
		ID:          id,
		RoomName:    roomName,
		Grade:       grade,
		Section:     section,
		TeacherID:   teacherID,
		Capacity:    capacity,
		Description: description,
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
