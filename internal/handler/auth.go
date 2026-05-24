package handler

import (
	"database/sql"
	"net/http"
	"time"

	"alhikmah-attendance-api/config"
	"alhikmah-attendance-api/internal/dto"
	"alhikmah-attendance-api/pkg/jwt"
	"alhikmah-attendance-api/pkg/response"
	"alhikmah-attendance-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	DB     *sql.DB
	Config config.Config
}

// DTOs moved to internal/dto

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body: "+err.Error()))
		return
	}

	var (
		id       string
		username string
		fullName string
		email    string
		role     string
		hash     string
	)

	err := h.DB.QueryRow(`
		SELECT id, username, full_name, email, role, password_hash
		FROM users 
		WHERE username = $1 AND is_active = true AND deleted_at IS NULL
	`, req.Username).Scan(&id, &username, &fullName, &email, &role, &hash)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, response.Error("Invalid credentials, Not Found"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.Error("Internal server error"))
		return
	}

	if !utils.CheckPasswordHash(req.Password, hash) {
		c.JSON(http.StatusUnauthorized, response.Error("Invalid credentials, hash mismatch"))
		return
	}

	// Update last login
	_, _ = h.DB.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1", id)

	token, err := jwt.GenerateToken(id, email, fullName, role, h.Config.JWTSecret, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to generate token"))
		return
	}

	refreshToken, err := jwt.GenerateToken(id, email, fullName, role, h.Config.JWTRefreshSecret, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to generate refresh token"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Login successful", gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":       id,
			"username": username,
			"name":     fullName,
			"email":    email,
			"role":     role,
		},
	}))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("Invalid request body: "+err.Error()))
		return
	}

	claims, err := jwt.ValidateToken(req.RefreshToken, h.Config.JWTRefreshSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error("Invalid refresh token"))
		return
	}

	// Generate new access token
	newToken, err := jwt.GenerateToken(claims.Sub, claims.Email, claims.Name, claims.Role, h.Config.JWTSecret, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to generate new token"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Token refreshed successfully", gin.H{
		"token": newToken,
	}))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// For JWT, logout is typically handled client-side by destroying the token.
	// You could implement a token blacklist here if needed.
	c.JSON(http.StatusOK, response.Success("Logout successful", nil))
}

func (h *AuthHandler) ResetPasswordRequest(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success("Password reset requested (stub)", nil))
}

func (h *AuthHandler) ResetPasswordConfirm(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success("Password reset confirmed (stub)", nil))
}

func (h *AuthHandler) ResetPasswordChange(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success("Password changed successfully (stub)", nil))
}
