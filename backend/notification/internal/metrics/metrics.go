package metrics

import (
	"net/http"

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

	NotificationKafkaEventReceivedSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notification_kafka_event_received_success_total",
			Help: "total number of notification kafka event received success",
		},
	)

	NotificationKafkaEventReceivedFailure = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notification_kafka_event_received_failure_total",
			Help: "total number of notification kafka event received failure",
		},
	)
)

func init() {
	prometheus.MustRegister(NotificationSuccess)
	prometheus.MustRegister(NotificationFailure)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		sw := &StatusResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(sw, r)
	})
}
