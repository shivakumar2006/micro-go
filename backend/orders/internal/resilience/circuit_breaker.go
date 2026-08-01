package resilience

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
	cb *gobreaker.CircuitBreaker
}

func NewCircuitBreaker() *CircuitBreaker {
	settings := gobreaker.Settings{
		Name: "cart-service",

		MaxRequests: 1,

		Interval: time.Minute,
		Timeout:  10 * time.Second,

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},

		OnStateChange: func(name string, from, to gobreaker.State) {
			slog.Info("circuit breaker state changed", slog.String("service", name), slog.String("from", from.String()), slog.String("to", to.String()))
		},
	}

	return &CircuitBreaker{
		cb: gobreaker.NewCircuitBreaker(settings),
	}
}

func Execute[T any](c *CircuitBreaker, fn func() (T, error)) (T, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		return fn()
	})

	if err != nil {
		var zero T
		return zero, err
	}

	typed, ok := result.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("unexpected response type")
	}

	return typed, nil
}
