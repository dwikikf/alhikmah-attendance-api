package dto

type CreateClassRequest struct {
	ClassName    string `json:"class_name" binding:"required,min=3"`
	TeacherID    string `json:"teacher_id" binding:"required,uuid"`
	AcademicYear string `json:"academic_year" binding:"required,len=9"`
	Capacity     int    `json:"capacity" binding:"required,gt=0,max=40"`
	Description  string `json:"description"`
}

type UpdateClassRequest struct {
	ClassName   string `json:"class_name" binding:"required,min=3"`
	TeacherID   string `json:"teacher_id" binding:"required,uuid"`
	Capacity    int    `json:"capacity" binding:"required,gt=0,max=40"`
	Description string `json:"description"`
}
