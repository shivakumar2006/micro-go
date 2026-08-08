package resilience

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
	Cb *gobreaker.CircuitBreaker
}

func NewCircuitBreaker() *CircuitBreaker {
	settings := gobreaker.Settings{
		Name: "order-service",

		// number of requests are allowed in half open state
		MaxRequests: 1,

		Interval: time.Minute,

		// open state lasts
		Timeout: 10 * time.Second,

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},

		OnStateChange: func(name string, from, to gobreaker.State) {
			slog.Info("circuit breaker state changed", slog.String("service", name), slog.String("from", from.String()), slog.String("to", to.String()))
		},
	}

	return &CircuitBreaker{
		Cb: gobreaker.NewCircuitBreaker(settings),
	}
}

func Execute[T any](c *CircuitBreaker, fn func() (T, error)) (T, error) {
	result, err := c.Cb.Execute(func() (interface{}, error) {
		return fn()
	})
	if err != nil {
		var zero T
		return zero, err
	}

	types, ok := result.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("unexpected response type")
	}

	return types, nil
}
