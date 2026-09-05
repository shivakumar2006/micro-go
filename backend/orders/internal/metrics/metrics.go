package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type StatusResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (s *StatusResponseWriter) WriteHeader(code int) {
	s.StatusCode = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *StatusResponseWriter) Write(b []byte) (int, error) {
	if s.StatusCode == 0 {
		s.StatusCode = http.StatusOK
	}
	return s.ResponseWriter.Write(b)
}

var (
	OrderRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_request_total",
			Help: "total order requests",
		},
		[]string{"method", "path", "status"},
	)

	OrderRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "order_request_duration_seconds",
			Help: "duration of order requests in seconds",
		},
		[]string{"method", "path", "status"},
	)

	OrderCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "order_created_total",
			Help: "total number of orders created",
		},
	)

	OrderStatusUpdateTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "order_status_update_total",
			Help: "total number of orders updated",
		},
	)
)

func init() {
	prometheus.MustRegister(OrderRequestTotal)
	prometheus.MustRegister(OrderRequestDuration)
	prometheus.MustRegister(OrderCreatedTotal)
	prometheus.MustRegister(OrderStatusUpdateTotal)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		sw := &StatusResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(sw, r)

		status := strconv.Itoa(sw.StatusCode)
		duration := time.Since(start).Seconds()

		OrderRequestTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()

		OrderRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
	})
}
