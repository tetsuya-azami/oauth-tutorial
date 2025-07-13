package logger

import (
	"log/slog"
	"os"
)

type MyLogger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
}

type myLogger struct {
	logger *slog.Logger
}

func (l *myLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *myLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *myLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// json形式のログを出力するためのLoggerを生成する関数
func NewMyLogger() MyLogger {
	return &myLogger{
		logger: slog.New(slog.NewJSONHandler(os.Stderr, nil)),
	}
}
