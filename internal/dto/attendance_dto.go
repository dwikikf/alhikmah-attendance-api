package dto

import "time"

// --- Request DTOs ---

type ScanQRRequest struct {
	NISN    string `json:"nisn" binding:"required,len=10"`
	Subject string `json:"subject"`
}

type ManualAttendanceRequest struct {
	StudentID string  `json:"student_id" binding:"required,uuid"`
	Status    string  `json:"status" binding:"required,oneof=hadir izin sakit tanpa_keterangan"`
	Subject   string  `json:"subject"`
	Notes     *string `json:"notes"`
}

type UpdateAttendanceRequest struct {
	Status string  `json:"status" binding:"required,oneof=hadir izin sakit tanpa_keterangan"`
	Reason string  `json:"reason"`
}

// --- Response DTOs ---

type AttendanceResponse struct {
	ID             string     `json:"id"`
	StudentID      string     `json:"student_id"`
	ClassID        string     `json:"class_id"`
	AttendanceDate time.Time  `json:"attendance_date"`
	Subject        *string    `json:"subject,omitempty"`
	Status         string     `json:"status"`
	RecordedBy     string     `json:"recorded_by"`
	RecordedAt     time.Time  `json:"recorded_at"`
	ScannedAt      *time.Time `json:"scanned_at,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	IsManual       bool       `json:"is_manual"`
}
