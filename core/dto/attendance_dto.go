package dto

type ScanQRRequest struct {
	NISN string `json:"nisn" binding:"required,len=10"`
}

type ManualAttendanceRequest struct {
	StudentID string `json:"student_id" binding:"required,uuid"`
	Status    string `json:"status" binding:"required,oneof=hadir izin sakit tanpa_keterangan"`
	Notes     *string `json:"notes"`
}
