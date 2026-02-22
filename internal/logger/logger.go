package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Fatal(s string, param2 string, err error) {
	panic("unimplemented")
}

func New(level string) (*Logger, error) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	logger := slog.New(handler)
	return &Logger{logger}, nil
}

func (l *Logger) Sync() error {
	return nil
}
