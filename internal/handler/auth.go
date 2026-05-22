package handler

import (
	"database/sql"
	"net/http"
	"time"

	"alhikmah-attendance-api/config"
	"alhikmah-attendance-api/pkg/jwt"
	"alhikmah-attendance-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	DB     *sql.DB
	Config config.Config
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials, Not Found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, hash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials, hash mismatch"})
		return
	}

	// Update last login
	_, _ = h.DB.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1", id)

	token, err := jwt.GenerateToken(id, email, fullName, role, h.Config.JWTSecret, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken, err := jwt.GenerateToken(id, email, fullName, role, h.Config.JWTRefreshSecret, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token":         token,
			"refresh_token": refreshToken,
			"user": gin.H{
				"id":        id,
				"username":  username,
				"name":      fullName,
				"email":     email,
				"role":      role,
			},
		},
		"message": "Login successful",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	claims, err := jwt.ValidateToken(req.RefreshToken, h.Config.JWTRefreshSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new access token
	newToken, err := jwt.GenerateToken(claims.Sub, claims.Email, claims.Name, claims.Role, h.Config.JWTSecret, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": newToken,
		},
		"message": "Token refreshed successfully",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// For JWT, logout is typically handled client-side by destroying the token.
	// You could implement a token blacklist here if needed.
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

func (h *AuthHandler) ResetPasswordRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset requested (stub)",
	})
}

func (h *AuthHandler) ResetPasswordConfirm(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password reset confirmed (stub)",
	})
}

func (h *AuthHandler) ResetPasswordChange(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully (stub)",
	})
}
