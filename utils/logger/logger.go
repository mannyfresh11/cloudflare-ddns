package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	Log *slog.Logger
}

func NewLogger(level slog.Leveler) *Logger {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	return &Logger{
		Log: log,
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.Log.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.Log.Error(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.Log.Warn(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.Log.Debug(msg, args...)
}
