package service

import (
	"errors"
	"fmt"
	"log/slog"

	"alhikmah-attendance-api/internal/domain"
)

type studentService struct {
	repo domain.StudentRepository
}

func NewStudentService(repo domain.StudentRepository) domain.StudentService {
	return &studentService{repo: repo}
}

func (s *studentService) GenerateQRCodeData(nisn, fullName, className string) string {
	// Simple format: NISN|FullName|ClassName
	// In a real scenario, this could be a JSON string or an encrypted payload
	return fmt.Sprintf("%s|%s|%s", nisn, fullName, className)
}

func (s *studentService) Create(student *domain.Student) error {
	if student.NISN == "" || student.FullName == "" || student.ClassID == "" {
		return errors.New("missing required fields")
	}

	// Generate QR code data automatically
	student.QRCodeData = s.GenerateQRCodeData(student.NISN, student.FullName, student.ClassName)

	return s.repo.Create(student)
}

func (s *studentService) CreateBulk(students []*domain.Student) error {
	for _, student := range students {
		if student.NISN == "" || student.FullName == "" || student.ClassID == "" {
			slog.Error("Bulk import failed due to missing required fields", slog.String("nisn", student.NISN), slog.String("name", student.FullName))
			return errors.New("missing required fields in one or more students")
		}
		// Generate QR code data automatically
		student.QRCodeData = s.GenerateQRCodeData(student.NISN, student.FullName, student.ClassName)
	}

	err := s.repo.CreateBulk(students)
	if err == nil {
		slog.Info("Student bulk import successful", slog.Int("count", len(students)))
	} else {
		slog.Error("Student bulk import failed at repo level", slog.Any("error", err))
	}
	return err
}

func (s *studentService) GetByID(id string) (*domain.Student, error) {
	return s.repo.GetByID(id)
}

func (s *studentService) GetByClassID(classID string) ([]*domain.Student, error) {
	return s.repo.GetByClassID(classID)
}

func (s *studentService) GetAll(teacherID string, isActive *bool, classID, search string, page, limit int) ([]*domain.Student, int, error) {
	return s.repo.GetAll(teacherID, isActive, classID, search, page, limit)
}

func (s *studentService) Update(student *domain.Student) error {
	if student.NISN != "" && student.FullName != "" && student.ClassName != "" {
		student.QRCodeData = s.GenerateQRCodeData(student.NISN, student.FullName, student.ClassName)
	}
	return s.repo.Update(student)
}

func (s *studentService) SoftDelete(id string) error {
	return s.repo.SoftDelete(id)
}
