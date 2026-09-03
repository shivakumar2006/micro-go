package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var paymentRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "payment_requests_total",
		Help: "total number of payment requests",
	},
	[]string{"method", "path", "status"},
)

var paymentRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "payment_request_duration_seconds",
		Help: "duration of HTTP payment requests in seconds",
	},
	[]string{"method", "path", "status"},
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

func init() {
	prometheus.MustRegister(paymentRequests)
	prometheus.MustRegister(paymentRequestDuration)
}

func metricsMiddleware(next http.Handler) http.Handler {
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

		paymentRequests.WithLabelValues(r.Method, r.URL.Path, status).Inc()

		duration := time.Since(start).Seconds()
		paymentRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
	})
}
