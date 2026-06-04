package mapper

import (
        "alhikmah-attendance-api/internal/domain"
        "alhikmah-attendance-api/internal/dto"
)

func ToUserDTO(user *domain.User) *dto.UserResponse {
        if user == nil {
                return nil
        }
        return &dto.UserResponse{
                ID:        user.ID,
                Username:  user.Username,
                Email:     user.Email,
                FullName:  user.FullName,
                Role:      user.Role,
                IsActive:  user.IsActive,
                CreatedAt: user.CreatedAt,
                UpdatedAt: user.UpdatedAt,
                LastLogin: user.LastLogin,
        }
}
