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
	NotificationRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notification_request_total",
			Help: "total number of notification request",
		},
		[]string{"method", "path", "status"},
	)

	NotificationRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "notification_request_duration_seconds",
			Help: "duration of notification request in seconds",
		},
		[]string{"method", "path", "status"},
	)

	NotificationSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notification_success_total",
			Help: "total number of notification success",
		},
	)

	NotificationFailure = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notification_failure_total",
			Help: "total number of notification failures",
		},
	)
)

func init() {
	prometheus.MustRegister(NotificationRequestTotal)
	prometheus.MustRegister(NotificationRequestDuration)
	prometheus.MustRegister(NotificationSuccess)
	prometheus.MustRegister(NotificationFailure)
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

		NotificationRequestTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		NotificationRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
	})
}
