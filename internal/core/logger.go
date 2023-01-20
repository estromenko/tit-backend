package core

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger

func NewLogger(conf *Config) *Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if conf.Debug {
		logger = logger.Output(zerolog.NewConsoleWriter())
	}

	return &logger
}
