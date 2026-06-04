package dto

import "time"

// --- Request DTOs ---

type CreateClassRequest struct {
	RoomName     string `json:"room_name" binding:"required,min=3"`
	Grade        int    `json:"grade" binding:"required,min=1,max=12"`
	Section      *int   `json:"section,omitempty" binding:"omitempty,gt=0"`
	TeacherID    string `json:"teacher_id" binding:"required,uuid"`
	AcademicYear string `json:"academic_year" binding:"required,len=9"`
	Capacity     int    `json:"capacity" binding:"required,gt=0,max=40"`
	Description  string `json:"description"`
}

type UpdateClassRequest struct {
	RoomName    string `json:"room_name" binding:"required,min=3"`
	Grade       int    `json:"grade" binding:"required,min=1,max=12"`
	Section     *int   `json:"section,omitempty" binding:"omitempty,gt=0"`
	TeacherID   string `json:"teacher_id" binding:"required,uuid"`
	Capacity    int    `json:"capacity" binding:"required,gt=0,max=40"`
	Description string `json:"description"`
}

// --- Response DTOs ---

type ClassResponse struct {
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
}
