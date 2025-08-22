package logger

import (
	"os"
	"path/filepath"
	"strconv"
	"webhook/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	// Enable unix time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Change caller format
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = log.With().Caller().Logger()

	// Set log mode
	logMode := config.Get().LogMode

	if logMode == "pretty" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	// Set log level
	logLevel := config.Get().LogLevel

	var zeroLogLevel zerolog.Level
	if err := zeroLogLevel.UnmarshalText([]byte(logLevel)); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal log level")
	}

	zerolog.SetGlobalLevel(zeroLogLevel)

	log.Info().Str("logLevel", logLevel).Msg("Logger successfully initialized")
}
