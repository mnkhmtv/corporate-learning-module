package logger

import (
	"log/slog"
	"os"
)

// NewLogger creates a new structured logger
func NewLogger(env string) *slog.Logger {
	var handler slog.Handler

	if env == "production" {
		// JSON format for production
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Pretty text format for development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}
