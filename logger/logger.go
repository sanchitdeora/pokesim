package logger

import (
	"log/slog"
	"os"
	"time"
)

// Initialize a global logger
func InitLogger() *slog.Logger {
	// Define the handler with JSON output
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // Includes file name and line number
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				// Format timestamp as ISO8601
				a.Value = slog.StringValue(time.Now().Format(time.RFC3339))
			}
			return a
		},
	})

	// Create and set the global logger
	logger := slog.New(handler)
	slog.SetDefault(logger) // Set as the default logger
	return logger
}
