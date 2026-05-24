package dto

type ScanQRRequest struct {
	QRCodeData string `json:"qr_code_data" binding:"required"`
}

type ManualAttendanceRequest struct {
	ClassID    string   `json:"class_id" binding:"required"`
	StudentIDs []string `json:"student_ids" binding:"required"`
	Status     string   `json:"status" binding:"required"`
	Notes      string   `json:"notes"`
}
