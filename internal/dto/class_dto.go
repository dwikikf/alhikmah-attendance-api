package dto

type CreateClassRequest struct {
	ClassName    string `json:"class_name" binding:"required"`
	TeacherID    string `json:"teacher_id" binding:"required"`
	AcademicYear string `json:"academic_year" binding:"required"`
	Capacity     int    `json:"capacity"`
	Description  string `json:"description"`
}

type UpdateClassRequest struct {
	ClassName   string `json:"class_name" binding:"required"`
	TeacherID   string `json:"teacher_id" binding:"required"`
	Capacity    int    `json:"capacity"`
	Description string `json:"description"`
}
