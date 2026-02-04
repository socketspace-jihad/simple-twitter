package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type logger struct {
	*zerolog.Logger
	*os.File
}

type loggerConfig func(*logger)

var (
	Logger *logger
)

func (l *logger) Close() {
	l.File.Close()
}

func NewLogger(configs ...loggerConfig) *logger {
	l := &logger{
		Logger: &zerolog.Logger{},
	}
	for _, config := range configs {
		if config != nil {
			config(l)
		}
	}

	return l
}

func WithFile(path string) loggerConfig {
	return func(l *logger) {
		if path == "" {
			path = "/var/log/simple-twitter.log"
		}
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("CRITICAL: could not open the log path")
		}

		multi := io.MultiWriter(file, os.Stdout)
		lg := zerolog.New(multi).With().Timestamp().Logger()
		l.Logger = &lg
		l.File = file
	}
}

func init() {
	lg := NewLogger(
		WithFile(os.Getenv("LOG_PATH")),
	)

	Logger = lg
}
