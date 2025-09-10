package logger

import (
	"io"
	"log/slog"
)

var logger *slog.Logger

func Init(level slog.Level, output io.Writer) {
	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewTextHandler(output, handlerOpts)
	logger = slog.New(handler)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, err error, args ...any) {
	args = append(args, slog.Any("error", err))
	logger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}
