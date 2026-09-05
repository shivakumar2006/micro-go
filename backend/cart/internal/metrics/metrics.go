package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type statusResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (s *statusResponseWriter) WriteHeader(code int) {
	s.StatusCode = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusResponseWriter) Write(b []byte) (int, error) {
	if s.StatusCode == 0 {
		s.StatusCode = http.StatusOK
	}
	return s.ResponseWriter.Write(b)
}

var (
	CartRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cart_request_total",
			Help: "total cart request",
		},
		[]string{"method", "path", "status"},
	)

	CartRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "cart_request_duration_seconds",
			Help: "duration of cart request per second",
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(CartRequestTotal)
	prometheus.MustRegister(CartRequestDuration)
}

func MerticsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		sw := &statusResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(sw, r)

		status := strconv.Itoa(sw.StatusCode)
		duration := time.Since(start).Seconds()

		CartRequestTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()

		CartRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
	})
}
