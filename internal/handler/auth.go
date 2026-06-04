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
		log.Printf("Warning: invalid token duration '%s', using default: %s", val, defaultDuration)
		return defaultDuration
	}
	return d
}

// DTOs moved to internal/dto

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Failed to parse login request", slog.String("error", err.Error()), slog.String("ip", c.ClientIP()))
		c.JSON(http.StatusBadRequest, response.ValidationError("Validation failed", utils.FormatValidationError(err)))
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
			c.JSON(http.StatusUnauthorized, response.Error("Invalid username or password"))
			return
		}
		slog.Error("Database error during login", slog.String("username", req.Username), slog.Any("error", err))
		slog.Error("Database error during login", slog.String("username", req.Username), slog.Any("error", err))
		HandleAppError(c, err, "Internal server error occurred during login")
		return
	}

	if !utils.CheckPasswordHash(req.Password, hash) {
		slog.Warn("Login failed: invalid password", slog.String("username", req.Username), slog.String("ip", c.ClientIP()))
		c.JSON(http.StatusUnauthorized, response.Error("Invalid username or password"))
		return
	}

	// Update last login
	_, _ = h.DB.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1", id)

	// Read token duration from configuration, with safe fallback
	accessTokenDuration := parseTokenDuration(h.Config.AccessTokenDuration, 1*time.Hour)
	refreshTokenDuration := parseTokenDuration(h.Config.RefreshTokenDuration, 7*24*time.Hour)

	token, err := jwt.GenerateToken(id, email, fullName, role, h.Config.JWTSecret, accessTokenDuration)
	if err != nil {
		HandleAppError(c, err, "Failed to generate access token")
		return
	}

	refreshToken, err := jwt.GenerateToken(id, email, fullName, role, h.Config.JWTRefreshSecret, refreshTokenDuration)
	if err != nil {
		HandleAppError(c, err, "Failed to generate refresh token")
		return
	}

	isProd := h.Config.AppEnv == "production" || h.Config.AppEnv == "prod"

	sameSiteMode := http.SameSiteLaxMode
	if isProd {
		sameSiteMode = http.SameSiteNoneMode
	}

	// Set refresh token in HttpOnly cookie
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

	// Read access token duration from configuration
	accessTokenDuration := parseTokenDuration(h.Config.AccessTokenDuration, 1*time.Hour)

	// Generate new access token
	newToken, err := jwt.GenerateToken(claims.Sub, claims.Email, claims.Name, claims.Role, h.Config.JWTSecret, accessTokenDuration)
	if err != nil {
		HandleAppError(c, err, "Failed to generate new token")
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
	// Feature temporarily disabled by request
	c.JSON(http.StatusForbidden, response.Error("Forgot password feature is disabled. Please contact the school admin to reset your password. Thank you."))
}

func (h *AuthHandler) ResetPasswordConfirm(c *gin.Context) {
	c.JSON(http.StatusForbidden, response.Error("Forgot password feature is disabled. Please contact the school admin. Thank you."))
}

func (h *AuthHandler) ResetPasswordChange(c *gin.Context) {
	c.JSON(http.StatusForbidden, response.Error("Forgot password feature is disabled. Please contact the school admin. Thank you."))
}
