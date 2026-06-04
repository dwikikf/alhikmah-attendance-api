package dto

import "time"

// --- Request DTOs ---

type AssignTeacherRequest struct {
	TeacherID string `json:"teacher_id" binding:"required,uuid"`
	Subject   string `json:"subject" binding:"required"`
}

// --- Response DTOs ---

type ClassTeacherResponse struct {
	ID           string    `json:"id"`
	TeacherID    string    `json:"teacher_id"`
	TeacherName  string    `json:"teacher_name"`
	ClassID      string    `json:"class_id"`
	ClassDisplay string    `json:"class_display"`
	AcademicYear string    `json:"academic_year"`
	Subject      string    `json:"subject"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
