package resilience

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type Retry struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	ShouldRetry  func(error) bool
}

func NewRetry(maxAttempts int, initialDelay, maxDelay time.Duration, shouldRetry func(error) bool) *Retry {
	return &Retry{
		MaxAttempts:  maxAttempts,
		InitialDelay: initialDelay,
		MaxDelay:     maxDelay,
		ShouldRetry:  shouldRetry,
	}
}

func DoRetry[T any](r *Retry, fn func() (T, error)) (T, error) {
	var zero T

	delay := r.InitialDelay

	for attempt := 1; attempt <= r.MaxAttempts; attempt++ {
		val, err := fn()
		if err == nil {
			return val, nil
		}
		slog.Warn("Retrying request", slog.Int("attempt", attempt), slog.Duration("delay", delay), slog.String("error", err.Error()))

		if !r.ShouldRetry(err) {
			return zero, err
		}

		if attempt == r.MaxAttempts {
			return zero, fmt.Errorf("maximum retry attempts (%d) reached: %w", r.MaxAttempts, err)
		}

		time.Sleep(delay)

		// exponential backoff
		delay = delay * 2

		if delay > r.MaxDelay {
			delay = r.MaxDelay
		}
	}

	return zero, errors.New("max retried attempts reached")
}
