package middleware

import (
	"net/http"
	"strings"

	"alhikmah-attendance-api/config"
	"alhikmah-attendance-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		// Handle different formats: "Bearer <token>", "bearer <token>", or just "<token>"
		if strings.HasPrefix(strings.ToLower(clientToken), "bearer ") {
			clientToken = clientToken[7:]
		}
		clientToken = strings.TrimSpace(clientToken)

		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Format of Authorization Token"})
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(clientToken, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("userID", claims.Sub)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found in token"})
			c.Abort()
			return
		}

		roleStr := userRole.(string)
		isAllowed := false
		for _, role := range requiredRoles {
			if role == roleStr {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied for this role"})
			c.Abort()
			return
		}

		c.Next()
	}
}
