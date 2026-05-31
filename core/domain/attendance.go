package domain

import "time"

type Attendance struct {
	ID             string    `json:"id"`
	StudentID      string    `json:"student_id"`
	ClassID        string    `json:"class_id"`
	AttendanceDate time.Time `json:"attendance_date"`
	Status         string    `json:"status"`
	RecordedBy     string    `json:"recorded_by"`
	RecordedAt     time.Time `json:"recorded_at"`
	ScannedAt      *time.Time `json:"scanned_at,omitempty"`
	Notes          *string   `json:"notes,omitempty"`
	IsManual       bool      `json:"is_manual"`
}

type AttendanceAudit struct {
	ID           string    `json:"id"`
	AttendanceID string    `json:"attendance_id"`
	OldStatus    *string   `json:"old_status,omitempty"`
	NewStatus    string    `json:"new_status"`
	ChangedBy    string    `json:"changed_by"`
	ChangedAt    time.Time `json:"changed_at"`
	Reason       *string   `json:"reason,omitempty"`
}

type AttendanceRepository interface {
	MarkAttendance(attendance *Attendance) error
	UpdateAttendance(attendance *Attendance, audit *AttendanceAudit) error
	GetByID(id string) (*Attendance, error)
	GetByClassAndDate(classID string, date time.Time) ([]*Attendance, error)
}

type AttendanceService interface {
	ProcessQRScan(nisn string, teacherID string, role string) error
	ProcessManualAttendance(studentID string, status string, notes *string, teacherID string, role string) error
	UpdateManual(attendanceID, status, reason, changedBy string) error
	GetClassAttendanceForToday(classID string) ([]*Attendance, error)
	GetByClassAndDate(classID string, date time.Time) ([]*Attendance, error)
}
