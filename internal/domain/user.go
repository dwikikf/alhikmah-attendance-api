package domain

import (
        "time"
)

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FullName     string     `json:"full_name"`
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	LastLogin    *time.Time `json:"last_login"`
}

type UserRepository interface {
	GetByID(id string) (*User, error)
	GetAll(role string, isActive *bool, offset, limit int) ([]*User, int, error)
	Create(user *User) error
	Update(user *User) error
	SoftDelete(id string) error
}

type UserService interface {
	GetByID(id string) (*User, error)
	GetAll(role string, isActive *bool, page, limit int) ([]*User, int, error)
	Create(username, email, password, fullName, role string) error
	Update(id, fullName, email, password string) error
	SoftDelete(id string) error
}
