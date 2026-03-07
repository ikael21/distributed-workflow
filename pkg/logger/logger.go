package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	ServiceName string
	Level       string
}

func NewLogger(cfg Config) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339

	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Logger(), nil
}
