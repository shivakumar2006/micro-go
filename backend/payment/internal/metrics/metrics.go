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
	PaymentRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_requests_total",
			Help: "total number of payment requests",
		},
		[]string{"method", "path", "status"},
	)

	PaymentRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "payment_request_duration_seconds",
			Help: "duration of HTTP payment requests in seconds",
		},
		[]string{"method", "path", "status"},
	)

	PaymentSuccessTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "payment_success_total",
			Help: "total number of successfull payments",
		},
	)

	PaymentFailureTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "payment_failure_total",
			Help: "total number of failure payments",
		},
	)

	OutboxPublishSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "outbox_publish_success_total",
			Help: "total number of successfull outbox publish",
		},
	)

	OutboxPublishFailure = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "outbox_publish_failure_total",
			Help: "total number of failure outbox publish",
		},
	)
)

func init() {
	prometheus.MustRegister(PaymentRequests)
	prometheus.MustRegister(PaymentRequestDuration)
	prometheus.MustRegister(PaymentSuccessTotal)
	prometheus.MustRegister(PaymentFailureTotal)
	prometheus.MustRegister(OutboxPublishSuccess)
	prometheus.MustRegister(OutboxPublishFailure)
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

		PaymentRequests.WithLabelValues(r.Method, r.URL.Path, status).Inc()

		duration := time.Since(start).Seconds()
		PaymentRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
	})
}
