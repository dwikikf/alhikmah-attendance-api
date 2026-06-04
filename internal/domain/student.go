package domain

import (
        "time"
)

type Student struct {
	ID         string     `json:"id"`
	NISN       string     `json:"nisn"`
	FullName   string     `json:"full_name"`
	ClassID    string     `json:"class_id"`
	ClassName  string     `json:"class_name,omitempty"`
	DOB        *string    `json:"date_of_birth,omitempty"`
	Gender     *string    `json:"gender,omitempty"`
	PhotoURL   *string    `json:"photo_url,omitempty"`
	QRCodeData string     `json:"qr_code_data"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`

	// Derived fields
	AttendanceToday *string    `json:"attendance_today,omitempty"`
	ScannedAt       *time.Time `json:"scanned_at,omitempty"`
}

type StudentRepository interface {
	Create(student *Student) error
	CreateBulk(students []*Student) error
	GetByID(id string) (*Student, error)
	GetByClassID(classID string) ([]*Student, error)
	GetByNISN(nisn string) (*Student, error)
	GetAll(teacherID string, isActive *bool, classID, search string, page, limit int) ([]*Student, int, error)
	Update(student *Student) error
	SoftDelete(id string) error
}

type StudentService interface {
	Create(nisn, fullName, classID string, dob, gender *string) error
	// CreateBulkFromRaw creates students from raw CSV-parsed data (nisn, fullName, classID, className, gender)
	CreateBulkFromRaw(entries []StudentRawEntry) error
	GetByID(id string) (*Student, error)
	GetByClassID(classID string) ([]*Student, error)
	GetAll(teacherID string, isActive *bool, classID, search string, page, limit int) ([]*Student, int, error)
	GenerateQRCodeData(nisn, fullName, className string) string
	Update(id, nisn, fullName, classID, className string, dob, gender *string, isActive bool) error
	SoftDelete(id string) error
}

// StudentRawEntry is used for bulk import from CSV
type StudentRawEntry struct {
	NISN      string
	FullName  string
	ClassID   string
	ClassName string
	Gender    *string
}
