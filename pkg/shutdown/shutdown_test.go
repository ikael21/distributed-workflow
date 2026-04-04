package shutdown

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

type testService struct {
	closeFn func(context.Context) error
}

func (s testService) Close(ctx context.Context) error {
	return s.closeFn(ctx)
}

func TestManagerRunPhasesRunsInPhaseOrder(t *testing.T) {
	t.Parallel()

	var (
		mu    sync.Mutex
		steps []int
	)

	record := func(phase int) func(context.Context) error {
		return func(context.Context) error {
			mu.Lock()
			defer mu.Unlock()
			steps = append(steps, phase)
			return nil
		}
	}

	manager := NewManager(time.Second, zerolog.Nop())
	manager.Register(1, testService{closeFn: record(1)})
	manager.Register(0, testService{closeFn: record(0)})
	manager.Register(2, testService{closeFn: record(2)})

	ok := manager.runPhases(context.Background())
	if !ok {
		t.Fatal("runPhases() = false, want true")
	}

	want := []int{0, 1, 2}
	if len(steps) != len(want) {
		t.Fatalf("len(steps) = %d, want %d", len(steps), len(want))
	}

	for i := range want {
		if steps[i] != want[i] {
			t.Fatalf("steps[%d] = %d, want %d", i, steps[i], want[i])
		}
	}
}

func TestManagerRunPhasesStopsOnContextDone(t *testing.T) {
	t.Parallel()

	started := make(chan struct{})
	release := make(chan struct{})

	manager := NewManager(time.Second, zerolog.Nop())
	manager.Register(0, testService{
		closeFn: func(ctx context.Context) error {
			close(started)
			select {
			case <-release:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	})
	manager.Register(1, testService{
		closeFn: func(context.Context) error {
			t.Fatal("phase after timeout should not run")
			return nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)
	go func() {
		done <- manager.runPhases(ctx)
	}()

	<-started
	cancel()

	select {
	case ok := <-done:
		if ok {
			t.Fatal("runPhases() = true, want false")
		}
	case <-time.After(time.Second):
		t.Fatal("runPhases() did not return after context cancellation")
	}

	close(release)
}

func TestManagerRunPhasesContinuesAfterServiceError(t *testing.T) {
	t.Parallel()

	ranSecondPhase := false

	manager := NewManager(time.Second, zerolog.Nop())
	manager.Register(0, testService{
		closeFn: func(context.Context) error {
			return errors.New("boom")
		},
	})
	manager.Register(1, testService{
		closeFn: func(context.Context) error {
			ranSecondPhase = true
			return nil
		},
	})

	ok := manager.runPhases(context.Background())
	if !ok {
		t.Fatal("runPhases() = false, want true")
	}

	if !ranSecondPhase {
		t.Fatal("second phase did not run after service error")
	}
}
