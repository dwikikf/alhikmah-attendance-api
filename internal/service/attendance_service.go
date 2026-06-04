package service

import (
	"errors"
	"log/slog"
	"time"

	"alhikmah-attendance-api/internal/domain"
	"alhikmah-attendance-api/pkg/cache"
)

type attendanceService struct {
	repo        domain.AttendanceRepository
	studentRepo domain.StudentRepository
	classRepo   domain.ClassRepository
	cache       *cache.Cache
}

func NewAttendanceService(
	repo domain.AttendanceRepository, 
	studentRepo domain.StudentRepository, 
	classRepo domain.ClassRepository,
	c *cache.Cache,
) domain.AttendanceService {
	return &attendanceService{
		repo:        repo,
		studentRepo: studentRepo,
		classRepo:   classRepo,
		cache:       c,
	}
}

func (s *attendanceService) ProcessQRScan(nisn string, teacherID string, role string, subject string) error {
	subjectKey := subject
	if subjectKey == "" {
		subjectKey = "reguler"
	}
	
	// Check if NISN exists in cache to prevent duplicate processing
	cacheKey := "qr_scan_" + nisn + "_" + subjectKey
	if _, found := s.cache.Get(cacheKey); found {
		slog.Warn("Duplicate QR scan prevented by cache", slog.String("nisn", nisn))
		return errors.New("student just recently scanned QR")
	}

	// Look up student details by NISN
	student, err := s.studentRepo.GetByNISN(nisn)
	if err != nil {
		return errors.New("student not found")
	}

	// Validate via IsTeacherResponsibleForStudent (skip if admin)
	if role != "admin" {
		isResponsible, err := s.classRepo.IsTeacherResponsibleForStudent(student.ID, teacherID)
		if err != nil {
			return err
		}
		if !isResponsible {
			slog.Warn("Unauthorized QR scan attempt", slog.String("teacher_id", teacherID), slog.String("student_id", student.ID))
			return errors.New("teacher does not have access to this student's class")
		}
	}

	// Check if attendance already recorded for today
	today := time.Now().Truncate(24 * time.Hour)
	attendances, err := s.repo.GetByClassAndDate(student.ClassID, today, subject)
	if err == nil {
		for _, a := range attendances {
			if a.StudentID == student.ID {
				// Save to cache if already scanned today to prevent future db hits
				s.cache.Set(cacheKey, true, 10*time.Second)
				slog.Warn("Student already scanned today", slog.String("student_id", student.ID))
				return errors.New("student has already attended today")
			}
		}
	}

	// Insert attendance record
	now := time.Now()
	var subjectPtr *string
	if subject != "" {
		subjectPtr = &subject
	}

	attendance := &domain.Attendance{
		StudentID:      student.ID,
		ClassID:        student.ClassID,
		AttendanceDate: today,
		Subject:        subjectPtr,
		Status:         "hadir", // Default status for QR
		RecordedBy:     teacherID,
		ScannedAt:      &now,
		IsManual:       false,
	}

	err = s.repo.MarkAttendance(attendance)
	if err != nil {
		return err
	}

	// Save NISN to cache with a short expiration (e.g., 10 seconds)
	s.cache.Set(cacheKey, true, 10*time.Second)

	slog.Info("QR Scan Attendance Successful", 
		slog.String("student_id", student.ID), 
		slog.String("class_id", student.ClassID),
		slog.String("teacher_id", teacherID))

	return nil
}

func (s *attendanceService) ProcessManualAttendance(studentID string, status string, notes *string, teacherID string, role string, subject string) error {
	if status != "hadir" && status != "izin" && status != "sakit" && status != "tanpa_keterangan" {
		return errors.New("invalid attendance status")
	}

	// Look up student to get ClassID
	student, err := s.studentRepo.GetByID(studentID)
	if err != nil {
		return errors.New("student not found")
	}
	classID := student.ClassID

	// Validate via IsTeacherResponsibleForStudent (skip if admin)
	if role != "admin" {
		isResponsible, err := s.classRepo.IsTeacherResponsibleForStudent(studentID, teacherID)
		if err != nil {
			return err
		}
		if !isResponsible {
			return errors.New("teacher does not have access to this student's class")
		}
	}

	today := time.Now().Truncate(24 * time.Hour)

	// Check if already recorded
	attendances, err := s.repo.GetByClassAndDate(classID, today, subject)
	if err == nil {
		for _, a := range attendances {
			if a.StudentID == studentID {
				return errors.New("student has already attended today")
			}
		}
	}

	var subjectPtr *string
	if subject != "" {
		subjectPtr = &subject
	}

	attendance := &domain.Attendance{
		StudentID:      studentID,
		ClassID:        classID,
		AttendanceDate: today,
		Subject:        subjectPtr,
		Status:         status,
		RecordedBy:     teacherID,
		IsManual:       true,
		Notes:          notes,
	}

	err = s.repo.MarkAttendance(attendance)
	if err == nil {
		slog.Info("Manual attendance recorded", 
			slog.String("student_id", studentID), 
			slog.String("status", status),
			slog.String("teacher_id", teacherID))
	}
	return err
}

func (s *attendanceService) UpdateManual(attendanceID, status, reason, changedBy string) error {
	a, err := s.repo.GetByID(attendanceID)
	if err != nil {
		return errors.New("attendance data not found")
	}
	
	attendance := &domain.Attendance{
		ID:         attendanceID,
		Status:     status,
		RecordedBy: changedBy,
	}

	oldStatus := a.Status
	
	audit := &domain.AttendanceAudit{
		AttendanceID: attendanceID,
		OldStatus:    &oldStatus,
		NewStatus:    status,
		ChangedBy:    changedBy,
	}
	
	if reason != "" {
		audit.Reason = &reason
	}

	err = s.repo.UpdateAttendance(attendance, audit)
	if err == nil {
		slog.Info("Manual attendance updated",
			slog.String("attendance_id", attendanceID),
			slog.String("old_status", oldStatus),
			slog.String("new_status", status),
			slog.String("changed_by", changedBy))
	}
	return err
}

func (s *attendanceService) GetClassAttendanceForToday(classID string, subject string) ([]*domain.Attendance, error) {
	today := time.Now().Truncate(24 * time.Hour)
	attendances, err := s.repo.GetByClassAndDate(classID, today, subject)
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

func (s *attendanceService) GetByClassAndDate(classID string, date time.Time, subject string) ([]*domain.Attendance, error) {
	attendances, err := s.repo.GetByClassAndDate(classID, date, subject)
	if err != nil {
		return nil, err
	}
	return attendances, nil
}
