package service

import (
	"errors"
	"fmt"

	"alhikmah-attendance-api/internal/domain"
)

type studentService struct {
	repo domain.StudentRepository
}

func NewStudentService(repo domain.StudentRepository) domain.StudentService {
	return &studentService{repo: repo}
}

func (s *studentService) GenerateQRCodeData(nisn, fullName, classID string) string {
	// Simple format: NISN|FullName|ClassID
	// In a real scenario, this could be a JSON string or an encrypted payload
	return fmt.Sprintf("%s|%s|%s", nisn, fullName, classID)
}

func (s *studentService) Create(student *domain.Student) error {
	if student.NISN == "" || student.FullName == "" || student.ClassID == "" {
		return errors.New("missing required fields")
	}

	// Generate QR code data automatically
	student.QRCodeData = s.GenerateQRCodeData(student.NISN, student.FullName, student.ClassID)

	return s.repo.Create(student)
}

func (s *studentService) GetByID(id string) (*domain.Student, error) {
	return s.repo.GetByID(id)
}

func (s *studentService) GetByClassID(classID string) ([]*domain.Student, error) {
	return s.repo.GetByClassID(classID)
}

func (s *studentService) GetAll(isActive *bool, page, limit int) ([]*domain.Student, int, error) {
	return s.repo.GetAll(isActive, page, limit)
}

func (s *studentService) Update(student *domain.Student) error {
	return s.repo.Update(student)
}

func (s *studentService) SoftDelete(id string) error {
	return s.repo.SoftDelete(id)
}
