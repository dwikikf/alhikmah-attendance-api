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
		role     string
		hash     string
	)

	err := h.DB.QueryRow(`
		SELECT id, username, full_name, role, password_hash
		FROM users 
		WHERE username = $1 AND is_active = true
	`, req.Username).Scan(&id, &username, &fullName, &role, &hash)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, hash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Update last login
	_, _ = h.DB.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1", id)

	token, err := jwt.GenerateToken(id, username, role, h.Config.JWTSecret, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":        id,
				"username":  username,
				"full_name": fullName,
				"role":      role,
			},
		},
		"message": "Login successful",
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
