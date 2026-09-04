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
	vehicleRequest = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vehicle_request_total",
			Help: "Total number of vehicle request",
		},
		[]string{"method", "path", "status"},
	)

	vehicleRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "vehicle_request_duration_seconds",
			Help: "Duration of vehicle request in seconds",
		},
		[]string{"method", "path", "status"},
	)

	vehicleCreateSuccessTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "vehicle_create_success_total",
			Help: "total number of vehicle create successfully",
		},
	)

	vehicleCreateFailureTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "vehicle_create_failure_total",
			Help: "total number of vehicle create failure",
		},
	)

	vehicleUpdateSuccessTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "vehicle_update_success_total",
			Help: "total number of vehicle update successfully",
		},
	)

	vehicleUpdateFailureTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "vehicle_update_failure_total",
			Help: "total number of vehicle update failure",
		},
	)

	vehicleDeleteSuccessTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "vehicle_delete_success_total",
			Help: "total number of vehicle delete successfully",
		},
	)

	vehicleDeleteFailureTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "vehicle_delete_failure_total",
			Help: "total number of vehicle delete failure",
		},
	)

	vehicleByStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vehicle_by_status",
			Help: "total number of vehicle by status",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(vehicleRequest)
	prometheus.MustRegister(vehicleRequestDuration)
	prometheus.MustRegister(vehicleCreateSuccessTotal)
	prometheus.MustRegister(vehicleCreateFailureTotal)
	prometheus.MustRegister(vehicleUpdateSuccessTotal)
	prometheus.MustRegister(vehicleUpdateFailureTotal)
	prometheus.MustRegister(vehicleDeleteSuccessTotal)
	prometheus.MustRegister(vehicleDeleteFailureTotal)
	prometheus.MustRegister(vehicleByStatus)
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

		status := strconv.Itoa(sw.StatusCode)
		vehicleRequest.WithLabelValues(r.Method, r.URL.Path, status).Inc()

		duration := time.Since(start).Seconds()
		vehicleRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
	})
}
