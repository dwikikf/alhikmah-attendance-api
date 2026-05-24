package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required,min=3"`
	Role     string `json:"role" binding:"required,oneof=admin teacher"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
}
