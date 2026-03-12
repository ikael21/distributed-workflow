package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	localEnv string = "local"
)

type Config struct {
	ServiceName string
	Level       string
	Env         string
}

func NewLogger(cfg Config) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339

	if cfg.Env == localEnv {
		return log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().
			Timestamp().
			Str("service", cfg.ServiceName).
			Logger(), nil
	}

	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Logger(), nil
}
