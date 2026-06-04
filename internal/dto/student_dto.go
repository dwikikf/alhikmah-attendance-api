package dto

import "time"

// --- Request DTOs ---

type CreateStudentRequest struct {
	NISN     string  `json:"nisn" binding:"required,len=10,numeric"`
	FullName string  `json:"full_name" binding:"required,min=3"`
	ClassID  string  `json:"class_id" binding:"required,uuid"`
	DOB      *string `json:"date_of_birth" binding:"omitempty,datetime=2006-01-02"`
	Gender   *string `json:"gender" binding:"required,oneof=laki-laki perempuan"`
}

type UpdateStudentRequest struct {
	FullName string  `json:"full_name" binding:"required,min=3"`
	ClassID  string  `json:"class_id" binding:"required,uuid"`
	DOB      *string `json:"date_of_birth" binding:"omitempty,datetime=2006-01-02"`
	Gender   *string `json:"gender" binding:"required,oneof=laki-laki perempuan"`
	IsActive *bool   `json:"is_active" binding:"required"`
}

// --- Response DTOs ---

type StudentResponse struct {
	ID              string     `json:"id"`
	NISN            string     `json:"nisn"`
	FullName        string     `json:"full_name"`
	ClassID         string     `json:"class_id"`
	ClassName       string     `json:"class_name,omitempty"`
	DOB             *string    `json:"date_of_birth,omitempty"`
	Gender          *string    `json:"gender,omitempty"`
	PhotoURL        *string    `json:"photo_url,omitempty"`
	QRCodeData      string     `json:"qr_code_data"`
	IsActive        bool       `json:"is_active"`
	AttendanceToday *string    `json:"attendance_today,omitempty"`
	ScannedAt       *time.Time `json:"scanned_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
