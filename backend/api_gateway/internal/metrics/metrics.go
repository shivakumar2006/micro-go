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
	GatewayRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_request_total",
			Help: "Total number of requests to the API Gateway",
		},
		[]string{"method", "path", "status"},
	)

	GatewayRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_request_duration_seconds",
			Help:    "Duration of HTTP requests handled by the API Gateway",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}, // 5ms to 10s
		},
		[]string{"method", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(GatewayRequestTotal)
	prometheus.MustRegister(GatewayRequestDuration)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &StatusResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(sw, r)

		status := strconv.Itoa(sw.StatusCode)
		duration := time.Since(start).Seconds()

		GatewayRequestTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		GatewayRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
	})
}
