package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	handler *slog.Logger
}

func New() *Logger {
	return &Logger{
		handler: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (l *Logger) Info(message string, args ...any) {
	l.handler.Info(message, args...)
}

func (l *Logger) Error(message string, args ...any) {
	l.handler.Error(message, args...)
}
