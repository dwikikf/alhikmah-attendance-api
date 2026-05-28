package handler

import (
	"database/sql"
	"log"
	"log/slog"
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

// parseTokenDuration membaca durasi dari string config, fallback ke nilai default jika tidak valid.
func parseTokenDuration(val string, defaultDuration time.Duration) time.Duration {
	if val == "" {
		return defaultDuration
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("Peringatan: nilai durasi token '%s' tidak valid, menggunakan default: %s", val, defaultDuration)
		return defaultDuration
	}
	return d
}

// DTOs moved to internal/dto

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Failed to parse login request", slog.String("error", err.Error()), slog.String("ip", c.ClientIP()))
		c.JSON(http.StatusBadRequest, response.ValidationError("Validasi gagal", utils.FormatValidationError(err)))
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
			slog.Warn("Login failed: user not found or inactive", slog.String("username", req.Username), slog.String("ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, response.Error("Invalid credentials, Not Found"))
			return
		}
		slog.Error("Database error during login", slog.String("username", req.Username), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, response.Error("Internal server error"))
		return
	}

	if !utils.CheckPasswordHash(req.Password, hash) {
		slog.Warn("Login failed: invalid password", slog.String("username", req.Username), slog.String("ip", c.ClientIP()))
		c.JSON(http.StatusUnauthorized, response.Error("Invalid credentials, hash mismatch"))
		return
	}

	// Update last login
	_, _ = h.DB.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1", id)

	// Baca durasi token dari konfigurasi, dengan fallback yang aman
	accessTokenDuration := parseTokenDuration(h.Config.AccessTokenDuration, 1*time.Hour)
	refreshTokenDuration := parseTokenDuration(h.Config.RefreshTokenDuration, 7*24*time.Hour)

	token, err := jwt.GenerateToken(id, email, fullName, role, h.Config.JWTSecret, accessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to generate token"))
		return
	}

	refreshToken, err := jwt.GenerateToken(id, email, fullName, role, h.Config.JWTRefreshSecret, refreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to generate refresh token"))
		return
	}

	isProd := h.Config.AppEnv == "production" || h.Config.AppEnv == "prod"

	sameSiteMode := http.SameSiteLaxMode
	if isProd {
		sameSiteMode = http.SameSiteNoneMode
	}

	// Set refresh token dalam HttpOnly cookie
	// Name, Value, MaxAge (detik), Path, Domain, Secure, HttpOnly
	c.SetSameSite(sameSiteMode)
	c.SetCookie("refresh_token", refreshToken, int(refreshTokenDuration/time.Second), "/", "", isProd, true)

	slog.Info("User logged in successfully", slog.String("username", username), slog.String("role", role), slog.String("ip", c.ClientIP()))

	c.JSON(http.StatusOK, response.Success("Login successful", gin.H{
		"token": token,
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
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, response.Error("Refresh token not found in cookie"))
		return
	}

	claims, err := jwt.ValidateToken(refreshToken, h.Config.JWTRefreshSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error("Invalid refresh token"))
		return
	}

	// Baca durasi access token dari konfigurasi
	accessTokenDuration := parseTokenDuration(h.Config.AccessTokenDuration, 1*time.Hour)

	// Generate new access token
	newToken, err := jwt.GenerateToken(claims.Sub, claims.Email, claims.Name, claims.Role, h.Config.JWTSecret, accessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("Failed to generate new token"))
		return
	}

	c.JSON(http.StatusOK, response.Success("Token refreshed successfully", gin.H{
		"token": newToken,
	}))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	isProd := h.Config.AppEnv == "production" || h.Config.AppEnv == "prod"

	sameSiteMode := http.SameSiteLaxMode
	if isProd {
		sameSiteMode = http.SameSiteNoneMode
	}

	c.SetSameSite(sameSiteMode)
	c.SetCookie("refresh_token", "", -1, "/", "", isProd, true)
	c.JSON(http.StatusOK, response.Success("Logout successful", nil))
}

func (h *AuthHandler) ResetPasswordRequest(c *gin.Context) {
	// Fitur dinonaktifkan sementara berdasarkan permintaan
	c.JSON(http.StatusForbidden, response.Error("Fitur lupa password dinonaktifkan. Silakan hubungi admin sekolah untuk mereset password Anda. Terima kasih."))
}

func (h *AuthHandler) ResetPasswordConfirm(c *gin.Context) {
	c.JSON(http.StatusForbidden, response.Error("Fitur lupa password dinonaktifkan. Silakan hubungi admin sekolah. Terima kasih."))
}

func (h *AuthHandler) ResetPasswordChange(c *gin.Context) {
	c.JSON(http.StatusForbidden, response.Error("Fitur lupa password dinonaktifkan. Silakan hubungi admin sekolah. Terima kasih."))
}
