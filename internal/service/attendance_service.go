package service

import (
	"errors"
	"strings"
	"time"

	"alhikmah-attendance-api/internal/domain"
)

type attendanceService struct {
	repo        domain.AttendanceRepository
	studentRepo domain.StudentRepository
}

func NewAttendanceService(repo domain.AttendanceRepository, studentRepo domain.StudentRepository) domain.AttendanceService {
	return &attendanceService{
		repo:        repo,
		studentRepo: studentRepo,
	}
}

func (s *attendanceService) ScanQR(qrData string, recordedBy string) error {
	// qrData format: NISN|FullName|ClassID
	parts := strings.Split(qrData, "|")
	if len(parts) < 3 {
		return errors.New("invalid QR code format")
	}
	nisn := parts[0]
	// classID := parts[2] // Could validate against this

	// 1. Validate student exists and is active
	student, err := s.studentRepo.GetByNISN(nisn)
	if err != nil {
		return errors.New("student not found or inactive")
	}

	// 2. Check if attendance already recorded for today
	today := time.Now().Truncate(24 * time.Hour)
	attendances, err := s.repo.GetByClassAndDate(student.ClassID, today)
	if err == nil {
		for _, a := range attendances {
			if a.StudentID == student.ID {
				return errors.New("student already scanned today")
			}
		}
	}

	// 3. Mark attendance
	now := time.Now()
	attendance := &domain.Attendance{
		StudentID:      student.ID,
		ClassID:        student.ClassID,
		AttendanceDate: today,
		Status:         "hadir",
		RecordedBy:     recordedBy,
		ScannedAt:      &now,
		IsManual:       false,
	}

	return s.repo.MarkAttendance(attendance)
}

func (s *attendanceService) ManualInput(classID string, studentIDs []string, status, notes, recordedBy string) error {
	if status != "hadir" && status != "izin" && status != "sakit" && status != "tanpa_keterangan" {
		return errors.New("invalid status")
	}

	today := time.Now().Truncate(24 * time.Hour)
	
	// In a real app, you would use a transaction to insert multiple records.
	// For simplicity, we loop here.
	for _, studentID := range studentIDs {
		// Ideally check if already exists to prevent duplicate key error
		attendance := &domain.Attendance{
			StudentID:      studentID,
			ClassID:        classID,
			AttendanceDate: today,
			Status:         status,
			RecordedBy:     recordedBy,
			IsManual:       true,
		}
		if notes != "" {
			attendance.Notes = &notes
		}
		
		// Ignore duplicate errors for bulk inserts (simplified)
		_ = s.repo.MarkAttendance(attendance)
	}

	return nil
}

func (s *attendanceService) UpdateManual(attendanceID, status, reason, changedBy string) error {
	a, err := s.repo.GetByID(attendanceID)
	if err != nil {
		return errors.New("attendance record not found")
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

	return s.repo.UpdateAttendance(attendance, audit)
}

func (s *attendanceService) GetClassAttendanceForToday(classID string) ([]*domain.Attendance, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return s.repo.GetByClassAndDate(classID, today)
}

func (s *attendanceService) GetByClassAndDate(classID string, date time.Time) ([]*domain.Attendance, error) {
	return s.repo.GetByClassAndDate(classID, date)
}
