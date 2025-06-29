package logger

import (
	"log/slog"
	"os"
)

// NewMyLogger initializes a new slog logger that outputs JSON formatted logs to stderr.
func NewMyLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil))
}
