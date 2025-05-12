package logger

import (
	"log/slog"
	"os"
	"sync"
)

const (
	envLocal = "local"
)

type Logger struct {
	*slog.Logger
}

var (
	logger Logger
	once   sync.Once
)

func Get() *Logger {
	once.Do(func() {
		var slogger *slog.Logger

		slogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

		logger = Logger{slogger}
	})
	return &logger
}
