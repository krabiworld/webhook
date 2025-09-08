package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"webhook/internal/config"
)

func Init() {
	logLevel := config.Get().LogLevel

	var slogLogLevel slog.Level
	if err := slogLogLevel.UnmarshalText([]byte(logLevel)); err != nil {
		slog.Error("Failed to unmarshal log level", "err", err.Error())
		os.Exit(1)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slogLogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				if source != nil {
					source.File = filepath.Base(source.File)
				}
			}
			return a
		},
	})

	slog.SetDefault(slog.New(handler))

	slog.Info("Logger initialized", "logLevel", slogLogLevel)
}
