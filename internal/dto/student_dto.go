package dto

type CreateStudentRequest struct {
	NISN     string  `json:"nisn" binding:"required,len=10"`
	FullName string  `json:"full_name" binding:"required,min=3"`
	ClassID  string  `json:"class_id" binding:"required,uuid"`
	Gender   *string `json:"gender"`
	DOB      *string `json:"date_of_birth"`
}

type UpdateStudentRequest struct {
	FullName string  `json:"full_name" binding:"required,min=3"`
	ClassID  string  `json:"class_id" binding:"required,uuid"`
	DOB      *string `json:"date_of_birth"`
	Gender   *string `json:"gender"`
	IsActive bool    `json:"is_active"`
}
