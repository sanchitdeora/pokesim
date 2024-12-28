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

func InitTestLogger() *slog.Logger {
	// Define the handler with JSON output
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,            // Includes file name and line number
		Level:     slog.LevelDebug, // Enable debug logs
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

type Logger interface {
	Log(message string)
}

type channelLogger struct {
	LogsChan chan<- string // Sends logs back to UI
}
func NewChannelLogger(logsChan chan<- string) Logger {
	return &channelLogger{LogsChan: logsChan}
}

type defaultLogger struct {}
func NewDefaultLogger() Logger {
	return &defaultLogger{}
}

func (c *channelLogger) Log(message string) {
	c.LogsChan <- message
}

func (c *defaultLogger) Log(message string) {
	slog.Info(message)
}
