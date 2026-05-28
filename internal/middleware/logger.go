package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs each HTTP request using slog
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		end := time.Now()
		latency := end.Sub(start)

		// Get request details
		method := c.Request.Method
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Prepare log attributes
		attrs := []any{
			slog.Int("status", statusCode),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("query", query),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
		}

		if errorMessage != "" {
			attrs = append(attrs, slog.String("error", errorMessage))
		}

		// Log based on status code
		if statusCode >= 500 {
			slog.Error("Server error", attrs...)
		} else if statusCode >= 400 {
			slog.Warn("Client error", attrs...)
		} else {
			slog.Info("Request completed", attrs...)
		}
	}
}
