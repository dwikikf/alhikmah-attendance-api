package domain

import "time"

type Class struct {
	ID           string    `json:"id"`
	ClassName    string    `json:"class_name"`
	TeacherID    string    `json:"teacher_id"`
	TeacherName  string    `json:"teacher_name"`
	AcademicYear string    `json:"academic_year"`
	Capacity     int       `json:"capacity"`
	StudentCount int       `json:"student_count"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ClassRepository interface {
	GetAll(teacherID string, academicYear string, offset, limit int) ([]*Class, int, error)
	GetByID(id string) (*Class, error)
}

type ClassService interface {
	GetAll(teacherID string, academicYear string, page, limit int) ([]*Class, int, error)
	GetByID(id string) (*Class, error)
}
