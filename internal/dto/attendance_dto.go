package dto

type ScanQRRequest struct {
	NISN string `json:"nisn" binding:"required"`
}

type ManualAttendanceRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
	Notes     string `json:"notes"`
}
