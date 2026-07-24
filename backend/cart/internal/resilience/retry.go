package resilience

import (
	"errors"
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

func RetryDo[T any](r *Retry, fn func() (T, error)) (T, error) {
	var zero T

	delay := r.InitialDelay

	for attempt := 1; attempt <= r.MaxAttempts; attempt++ {

		result, err := fn()
		if err == nil {
			return result, nil
		}

		// not retry
		if !r.ShouldRetry(err) {
			return zero, err
		}

		// last attempt
		if attempt == r.MaxAttempts {
			return zero, err
		}

		time.Sleep(delay)

		// exponetial backoff
		delay = delay * 2

		if delay > r.MaxDelay {
			delay = r.MaxDelay
		}
	}

	return zero, errors.New("retry failed")
}
