package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

// GracefulShutdown is implemented by components that release resources when the
// process is stopping (HTTP servers, connection pools, workers, etc.).
type GracefulShutdown interface {
	Close(context.Context) error
}

// Manager runs registered services in ordered phases after SIGINT/SIGTERM.
// Lower phase numbers complete first; within a phase, services close in parallel.
// A typical layout is phase 0 for HTTP (or RPC) and phase 1 for databases.
type Manager struct {
	logger  zerolog.Logger
	phases  map[int][]GracefulShutdown
	timeout time.Duration
}

// NewManager builds a shutdown coordinator. timeout bounds the entire shutdown
// sequence (all phases); a single context is passed to every Close call.
func NewManager(timeout time.Duration, logger zerolog.Logger) *Manager {
	return &Manager{
		logger:  logger,
		phases:  make(map[int][]GracefulShutdown),
		timeout: timeout,
	}
}

// Register adds a service to a phase. Phases run in ascending order (0, 1, 2…).
// Multiple Register calls with the same phase append to that phase’s group.
func (m *Manager) Register(phase int, gs GracefulShutdown) {
	m.phases[phase] = append(m.phases[phase], gs)
}

// Wait blocks until a shutdown signal, then runs phases in order until the
// global timeout elapses or every Close returns.
func (m *Manager) Wait() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	<-quit
	m.logger.Info().Msg("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	if m.runPhases(ctx) {
		m.logger.Info().Msg("all shutdown phases complete")
	}
}

func (m *Manager) runPhases(ctx context.Context) bool {
	keys := make([]int, 0, len(m.phases))
	for k := range m.phases {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, phase := range keys {
		services := m.phases[phase]
		m.logger.Info().Int("phase", phase).Msg("shutdown phase started")

		var wg sync.WaitGroup
		wg.Add(len(services))
		for _, service := range services {
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
			m.logger.Info().Int("phase", phase).Msg("phase complete")
		case <-ctx.Done():
			m.logger.Warn().Err(ctx.Err()).Int("phase", phase).Msg("shutdown timeout during phase")
			return false
		}
	}

	return true
}
