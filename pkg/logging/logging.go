package logging

import (
	"log/slog"
	"os"
)

func InitLog() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return log
}
