package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AppName       string `envconfig:"APP_NAME"    required:"true"`
	AppVersion    string `envconfig:"APP_VERSION" required:"true"`
	Level         string `default:"error"         envconfig:"LOGGER_LEVEL"`
	PrettyConsole bool   `default:"false"         envconfig:"LOGGER_PRETTY_CONSOLE"`
}

func Init(c Config) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	level, err := zerolog.ParseLevel(c.Level)
	if err == nil {
		zerolog.SetGlobalLevel(level)
	} else {
		log.Info().Msgf("Invalid log level %s specified", c.Level)
	}

	log.Logger = log.With().Caller().Str("app", c.AppName).Str("vsn", c.AppVersion).Logger()

	if c.PrettyConsole {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:01:02"})
	}

	log.Info().Msg("Logger initialized")
}
