package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusResponseWriter) WriteHeader(code int) {
	s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusResponseWriter) Write(b []byte) (int, error) {
	if s.statusCode == 0 {
		s.statusCode = http.StatusOK
	}

	return s.ResponseWriter.Write(b)
}

var (
	AuthRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_request_total",
			Help: "total number of auth request",
		},
		[]string{"method", "path", "status"},
	)

	AuthRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "auth_request_duration_seconds",
			Help: "duration of auth request in seconds",
		},
		[]string{"method", "path", "status"},
	)

	AuthLoginSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_login_success_total",
			Help: "total number of login success",
		},
		[]string{"method", "path", "status"},
	)

	AuthLoginFailureTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_login_failure_total",
			Help: "total number of login failure",
		},
		[]string{"method", "path", "status"},
	)

	AuthRefreshSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_refresh_success_total",
			Help: "total number of refresh success",
		},
		[]string{"method", "path", "status"},
	)

	AuthRefreshFailureTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_refresh_failure_total",
			Help: "total number of refresh failure",
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(AuthRequest)
	prometheus.MustRegister(AuthRequestDuration)
	prometheus.MustRegister(AuthLoginSuccessTotal)
	prometheus.MustRegister(AuthLoginFailureTotal)
	prometheus.MustRegister(AuthRefreshSuccessTotal)
	prometheus.MustRegister(AuthRefreshFailureTotal)
}

func MetricsMiddleware(next http.Handler) http.Handler {
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

		status := strconv.Itoa(sw.statusCode)
		duration := time.Since(start)

		AuthRequest.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		AuthRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration.Seconds())

		if sw.statusCode == http.StatusOK {
			AuthLoginSuccessTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		} else {
			AuthLoginFailureTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
		}
	})
}
