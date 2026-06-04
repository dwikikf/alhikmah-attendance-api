package service

import (
	"errors"
	"strings"

	"alhikmah-attendance-api/internal/domain"
)

type classTeacherService struct {
	repo      domain.ClassTeacherRepository
	userRepo  domain.UserRepository
	classRepo domain.ClassRepository
}

func NewClassTeacherService(repo domain.ClassTeacherRepository, userRepo domain.UserRepository, classRepo domain.ClassRepository) domain.ClassTeacherService {
	return &classTeacherService{
		repo:      repo,
		userRepo:  userRepo,
		classRepo: classRepo,
	}
}

func (s *classTeacherService) Assign(teacherID, classID, subject string) error {
	if strings.TrimSpace(subject) == "" {
		return errors.New("subject cannot be empty")
	}

	// Validate teacher exists
	user, err := s.userRepo.GetByID(teacherID)
	if err != nil || user == nil {
		return errors.New("teacher not found")
	}
	if user.Role != "teacher" {
		return errors.New("user is not a teacher")
	}

	// Validate class exists
	class, err := s.classRepo.GetByID(classID)
	if err != nil || class == nil {
		return errors.New("class not found")
	}

	ct := &domain.ClassTeacher{
		TeacherID:    teacherID,
		ClassID:      classID,
		Subject:      subject,
		AcademicYear: class.AcademicYear,
		Role:         "subject_teacher",
	}

	return s.repo.Assign(ct)
}

func (s *classTeacherService) Unassign(teacherID, classID, subject string) error {
	if strings.TrimSpace(subject) == "" {
		return errors.New("subject cannot be empty")
	}
	return s.repo.Unassign(teacherID, classID, subject)
}

func (s *classTeacherService) GetByClassID(classID string) ([]*domain.ClassTeacher, error) {
	teachers, err := s.repo.GetByClassID(classID)
	if err != nil {
		return nil, err
	}
	return teachers, nil
}

func (s *classTeacherService) GetSubjectAssignments(teacherID string) ([]*domain.ClassTeacher, error) {
	assignments, err := s.repo.GetSubjectAssignments(teacherID)
	if err != nil {
		return nil, err
	}
	return assignments, nil
}
