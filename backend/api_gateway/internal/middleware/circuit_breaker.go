package middleware

import (
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreakerConfig struct {
	FailureThreshold uint32
	MaxRequests      uint32
	Interval         time.Duration
	Timeout          time.Duration
}

type CircuitBreaker struct {
	mu       sync.RWMutex
	breakers map[string]*gobreaker.CircuitBreaker
	config   CircuitBreakerConfig
}

func NewCircuitBreakerConfig(cfg CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		breakers: make(map[string]*gobreaker.CircuitBreaker),
		config:   cfg,
	}
}

func (c *CircuitBreaker) GetBreaker(serviceName string, cfg CircuitBreakerConfig) *gobreaker.CircuitBreaker {
	// fast path (read lock)
	c.mu.RLock()
	breaker, exists := c.breakers[serviceName]
	c.mu.RUnlock()

	if exists {
		return breaker
	}

	// slow path (write lock)
	c.mu.Lock()
	defer c.mu.Unlock()

	// double check
	if breaker, exists := c.breakers[serviceName]; exists {
		return breaker
	}

	settings := gobreaker.Settings{
		Name:        serviceName,
		MaxRequests: c.config.MaxRequests,
		Interval:    c.config.Interval,
		Timeout:     c.config.Timeout,

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > c.config.FailureThreshold
		},

		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fromStr := stateToString(from)
			toStr := stateToString(to)
			println("Circuit Breaker [" + name + "]" + fromStr + " -> " + toStr)
		},
	}

	breaker = gobreaker.NewCircuitBreaker(settings)
	c.breakers[serviceName] = breaker
	return breaker
}

func stateToString(state gobreaker.State) string {
	switch state {
	case gobreaker.StateClosed:
		return "CLOSED"
	case gobreaker.StateHalfOpen:
		return "HALF_OPEN"
	case gobreaker.StateOpen:
		return "OPEN"
	default:
		return "UNKNOWN"
	}
}
