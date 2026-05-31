package dto

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
