package dto

import "time"

// --- Request DTOs ---

type EnrollRequest struct {
	StudentID    string `json:"student_id" binding:"required,uuid"`
	AcademicYear string `json:"academic_year" binding:"required,len=9"`
}

type PromoteItemRequest struct {
	StudentID     string `json:"student_id" binding:"required,uuid"`
	TargetClassID string `json:"target_class_id" binding:"required,uuid"`
}

type PromoteClassRequest struct {
	AcademicYear string               `json:"academic_year" binding:"required,len=9"`
	Items        []PromoteItemRequest `json:"items" binding:"required,gt=0"`
}

type TransferStudentRequest struct {
	FromClassID  string `json:"from_class_id" binding:"required,uuid"`
	ToClassID    string `json:"to_class_id" binding:"required,uuid"`
	AcademicYear string `json:"academic_year" binding:"required,len=9"`
}

// --- Response DTOs ---

type EnrollmentResponse struct {
	ID           string     `json:"id"`
	StudentID    string     `json:"student_id"`
	StudentName  string     `json:"student_name"`
	ClassID      string     `json:"class_id"`
	ClassDisplay string     `json:"class_display"`
	AcademicYear string     `json:"academic_year"`
	EnrolledAt   time.Time  `json:"enrolled_at"`
	EndedAt      *time.Time `json:"ended_at,omitempty"`
	EndReason    *string    `json:"end_reason,omitempty"`
}
type EnrollmentHistoryResponse = EnrollmentResponse
