package middleware

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

var errBackendFailure = errors.New("Backend failure")

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

func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		breakers: make(map[string]*gobreaker.CircuitBreaker),
		config:   cfg,
	}
}

func (c *CircuitBreaker) GetBreaker(serviceName string) *gobreaker.CircuitBreaker {
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

func (c *CircuitBreaker) Protect(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			breaker := c.GetBreaker(serviceName)
			response := NewResponseRecorder(w)

			_, err := breaker.Execute(func() (interface{}, error) {
				next.ServeHTTP(response, r)

				if response.StatusCode >= 500 {
					return nil, errBackendFailure
				}

				return nil, nil
			})

			if err == nil {
				return
			}

			switch {
			case errors.Is(err, gobreaker.ErrOpenState), errors.Is(err, gobreaker.ErrTooManyRequests):
				http.Error(w, "Service temporary unavailable", http.StatusServiceUnavailable)

			case errors.Is(err, errBackendFailure):
				return

			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		})
	}
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
	}
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
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
