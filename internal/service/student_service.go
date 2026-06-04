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
	return fmt.Sprintf("%s|%s|%s", nisn, fullName, className)
}

func (s *studentService) Create(nisn, fullName, classID string, dob, gender *string) error {
	if nisn == "" || fullName == "" || classID == "" {
		return errors.New("missing required fields")
	}

	student := &domain.Student{
		NISN:     nisn,
		FullName: fullName,
		ClassID:  classID,
		DOB:      dob,
		Gender:   gender,
	}
	// QR code data will be set by the repository after class name is resolved,
	// but we generate a basic one here; className will be empty if not provided.
	student.QRCodeData = s.GenerateQRCodeData(nisn, fullName, "")

	return s.repo.Create(student)
}

func (s *studentService) CreateBulkFromRaw(entries []domain.StudentRawEntry) error {
	if len(entries) == 0 {
		return errors.New("missing required fields in one or more students")
	}

	var students []*domain.Student
	for _, e := range entries {
		if e.NISN == "" || e.FullName == "" || e.ClassID == "" {
			slog.Error("Bulk import failed due to missing required fields", slog.String("nisn", e.NISN), slog.String("name", e.FullName))
			return errors.New("missing required fields in one or more students")
		}
		student := &domain.Student{
			NISN:      e.NISN,
			FullName:  e.FullName,
			ClassID:   e.ClassID,
			ClassName: e.ClassName,
			Gender:    e.Gender,
		}
		student.QRCodeData = s.GenerateQRCodeData(e.NISN, e.FullName, e.ClassName)
		students = append(students, student)
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
	student, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (s *studentService) GetByClassID(classID string) ([]*domain.Student, error) {
	students, err := s.repo.GetByClassID(classID)
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (s *studentService) GetAll(teacherID string, isActive *bool, classID, search string, page, limit int) ([]*domain.Student, int, error) {
	students, total, err := s.repo.GetAll(teacherID, isActive, classID, search, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return students, total, nil
}

func (s *studentService) Update(id, nisn, fullName, classID, className string, dob, gender *string, isActive bool) error {
	student := &domain.Student{
		ID:        id,
		NISN:      nisn,
		FullName:  fullName,
		ClassID:   classID,
		ClassName: className,
		DOB:       dob,
		Gender:    gender,
		IsActive:  isActive,
	}
	if nisn != "" && fullName != "" && className != "" {
		student.QRCodeData = s.GenerateQRCodeData(nisn, fullName, className)
	}
	return s.repo.Update(student)
}

func (s *studentService) SoftDelete(id string) error {
	return s.repo.SoftDelete(id)
}
