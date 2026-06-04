package domain

import (
        "time"
)

type StudentEnrollment struct {
	ID           string
	StudentID    string
	StudentName  string  // join dari students
	ClassID      string
	ClassDisplay string  // display name kelas
	AcademicYear string
	EnrolledAt   time.Time
	EndedAt      *time.Time
	EndReason    *string  // 'promoted' | 'transferred' | 'graduated' | 'dropped'
}

type PromoteItem struct {
	StudentID     string  // siswa yang akan dipromote
	TargetClassID string  // kelas tujuan (bisa beda section)
}

type EnrollmentRepository interface {
  Enroll(e *StudentEnrollment) error
	GetActiveByStudentID(studentID string) (*StudentEnrollment, error)
  GetActiveByClassID(classID string) ([]*StudentEnrollment, error)
  GetHistoryByStudentID(studentID string) ([]*StudentEnrollment, error)
  EndEnrollment(studentID, classID string, reason string) error
  BulkEnroll(items []PromoteItem, academicYear string) (int, error)
}

type EnrollmentService interface {
	Enroll(studentID, classID, academicYear string) error
	PromoteClass(items []PromoteItem, academicYear string) (int, error)
	TransferStudent(studentID, fromClassID, toClassID, academicYear string) error
	GetActiveByStudentID(studentID string) (*StudentEnrollment, error)
	GetActiveByClassID(classID string) ([]*StudentEnrollment, error)
	GetHistoryByStudentID(studentID string) ([]*StudentEnrollment, error)
}
