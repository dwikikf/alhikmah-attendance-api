package domain

import (
	"fmt"
	"time"
)

type Class struct {
	ID           string     `json:"id"`
	RoomName     string     `json:"room_name"`
	Grade        int        `json:"grade"`
	Section      *int       `json:"section"`
	DisplayName  string     `json:"display_name"`
	TeacherID    string     `json:"teacher_id"`
	TeacherName  string     `json:"teacher_name"`
	AcademicYear string     `json:"academic_year"`
	Capacity     int        `json:"capacity"`
	StudentCount int        `json:"student_count"`
	Description  string     `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

func (c *Class) GetDisplayName() string {
	if c.Section != nil {
		return fmt.Sprintf("Kelas %d %s %d", c.Grade, c.RoomName, *c.Section)
	}
	return fmt.Sprintf("Kelas %d %s", c.Grade, c.RoomName)
}

type ClassRepository interface {
	GetAll(teacherID string, academicYear string, offset, limit int) ([]*Class, int, error)
	GetByID(id string) (*Class, error)
	Create(class *Class) error
	Update(class *Class) error
	SoftDelete(id string) error
	IsTeacherResponsibleForStudent(studentID string, teacherID string) (bool, error)
}

type ClassService interface {
	GetAll(teacherID string, academicYear string, page, limit int) ([]*Class, int, error)
	GetByID(id string) (*Class, error)
	Create(roomName string, grade int, section *int, teacherID, academicYear string, capacity int, description string) error
	Update(id, roomName string, grade int, section *int, teacherID string, capacity int, description string) error
	SoftDelete(id string) error
}
