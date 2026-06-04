package service

import (
	"errors"
	"fmt"

	"alhikmah-attendance-api/internal/domain"
)

type enrollmentService struct {
	repo         domain.EnrollmentRepository
	studentRepo  domain.StudentRepository
	classRepo    domain.ClassRepository
}

func NewEnrollmentService(repo domain.EnrollmentRepository, studentRepo domain.StudentRepository, classRepo domain.ClassRepository) domain.EnrollmentService {
	return &enrollmentService{
		repo:        repo,
		studentRepo: studentRepo,
		classRepo:   classRepo,
	}
}

func (s *enrollmentService) Enroll(studentID, classID, academicYear string) error {
	// Cek siswa
	student, err := s.studentRepo.GetByID(studentID)
	if err != nil || student == nil {
		return errors.New("student not found")
	}

	// Cek kelas
	class, err := s.classRepo.GetByID(classID)
	if err != nil || class == nil {
		return errors.New("class not found")
	}

	// Cek apakah siswa sudah punya enrollment aktif di TA ini
	active, err := s.repo.GetActiveByStudentID(studentID)
	if err != nil {
		return err
	}
	if active != nil && active.AcademicYear == academicYear {
		return fmt.Errorf("student already has active enrollment in class %s for academic year %s", active.ClassDisplay, academicYear)
	}

	enrollment := &domain.StudentEnrollment{
		StudentID:    studentID,
		ClassID:      classID,
		AcademicYear: academicYear,
	}

	return s.repo.Enroll(enrollment)
}

func (s *enrollmentService) PromoteClass(items []domain.PromoteItem, academicYear string) (int, error) {
	if academicYear == "" || len(academicYear) != 9 || academicYear[4] != '/' {
		return 0, errors.New("academic_year is required and must be in format YYYY/YYYY")
	}

	if len(items) == 0 {
		return 0, errors.New("no students to promote")
	}

	return s.repo.BulkEnroll(items, academicYear)
}

func (s *enrollmentService) TransferStudent(studentID, fromClassID, toClassID, academicYear string) error {
	// Cek enrollment aktif
	active, err := s.repo.GetActiveByStudentID(studentID)
	if err != nil {
		return err
	}
	if active == nil || active.ClassID != fromClassID {
		return errors.New("student does not have an active enrollment in the specified from_class_id")
	}

	// Tutup enrollment lama
	err = s.repo.EndEnrollment(studentID, fromClassID, "transferred")
	if err != nil {
		return err
	}

	// Buka enrollment baru
	return s.Enroll(studentID, toClassID, academicYear)
}

func (s *enrollmentService) GetActiveByStudentID(studentID string) (*domain.StudentEnrollment, error) {
	enrollment, err := s.repo.GetActiveByStudentID(studentID)
	if err != nil {
		return nil, err
	}
	return enrollment, nil
}

func (s *enrollmentService) GetActiveByClassID(classID string) ([]*domain.StudentEnrollment, error) {
	enrollments, err := s.repo.GetActiveByClassID(classID)
	if err != nil {
		return nil, err
	}
	return enrollments, nil
}

func (s *enrollmentService) GetHistoryByStudentID(studentID string) ([]*domain.StudentEnrollment, error) {
	enrollments, err := s.repo.GetHistoryByStudentID(studentID)
	if err != nil {
		return nil, err
	}
	return enrollments, nil
}
