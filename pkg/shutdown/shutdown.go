package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type GracefulShutdown interface {
	Close(ctx context.Context) error
}

type manager struct {
	logger   zerolog.Logger
	services []GracefulShutdown
	timeout  time.Duration
}

func NewManager(timeout time.Duration, logger zerolog.Logger) *manager {
	return &manager{
		services: []GracefulShutdown{},
		timeout: timeout,
		logger: logger,
	}
}

func (m *manager) Register(gs GracefulShutdown) {
	m.services = append(m.services, gs)
}

func (m *manager) Wait() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit
	m.logger.Info().Msg("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, service := range m.services {
		wg.Add(1)
		go func(gs GracefulShutdown) {
			defer wg.Done()

			gsLogger := m.logger.With().Str("component", fmt.Sprintf("%T", gs)).Logger()
			if err := gs.Close(ctx); err != nil {
				gsLogger.Error().Err(err).Msg("service shutdown failed")
				return
			}
			gsLogger.Info().Msg("stopped service")
		}(service)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		m.logger.Info().Msg("all services gracefully shutdown")
	case <-ctx.Done():
		m.logger.Warn().Err(ctx.Err()).Msg("shutdown timeout reached")
	}
}
