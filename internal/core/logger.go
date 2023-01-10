package core

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(conf *Config) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if conf.Debug {
		logger = logger.Output(zerolog.NewConsoleWriter())
	}

	return &logger
}
