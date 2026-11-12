package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Setup настраивает глобальный slog: JSON в production, текст локально.
func Setup(env string) {
	level := slog.LevelInfo
	if strings.EqualFold(env, "local") || strings.EqualFold(env, "dev") {
		level = slog.LevelDebug
	}

	options := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if strings.EqualFold(env, "production") || strings.EqualFold(env, "prod") {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	slog.SetDefault(slog.New(handler))
}
