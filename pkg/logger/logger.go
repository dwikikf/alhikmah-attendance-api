package logger

import (
	"log/slog"
	"os"
)

// InitLogger initializes the global slog based on the environment.
// In development, it uses a text handler with DEBUG level.
// In production, it uses a JSON handler with INFO level.
func InitLogger(env string) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo, // Default for production
	}

	if env == "development" {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		// Production defaults to JSON format for Log Management Services (like Better Stack)
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	
	// Set this logger as the global default logger
	slog.SetDefault(logger)
}
