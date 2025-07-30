package mylogger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
}

type logger struct {
	logger *slog.Logger
}

func (l *logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// json形式のログを出力するためのLoggerを生成する関数
func NewLogger() Logger {
	return &logger{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}
