package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// New creates a new zerolog logger instance for a service
func New(serviceName string) zerolog.Logger {
	// Use console writer for pretty logging in development
	// In production, you'd remove this and just log JSON.
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}

	// Create a logger with a "service" field
	logger := zerolog.New(consoleWriter).With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	return logger
}
