package domain

import (
        "time"
)

type ClassTeacher struct {
	ID           string    `json:"id"`
	TeacherID    string    `json:"teacher_id"`
	TeacherName  string    `json:"teacher_name"`
	ClassID      string    `json:"class_id"`
	ClassDisplay string    `json:"class_display"`
	AcademicYear string    `json:"academic_year"`
	Subject      string    `json:"subject"`
	Role         string    `json:"role"` // 'homeroom' | 'subject_teacher'
	CreatedAt    time.Time `json:"created_at"`
}

type ClassTeacherRepository interface {
	Assign(ct *ClassTeacher) error
	Unassign(teacherID, classID, subject string) error
	GetByClassID(classID string) ([]*ClassTeacher, error)
	GetByTeacherID(teacherID string) ([]*ClassTeacher, error)
	GetSubjectAssignments(teacherID string) ([]*ClassTeacher, error)
}

type ClassTeacherService interface {
	Assign(teacherID, classID, subject string) error
	Unassign(teacherID, classID, subject string) error
	GetByClassID(classID string) ([]*ClassTeacher, error)
	GetSubjectAssignments(teacherID string) ([]*ClassTeacher, error)
}
