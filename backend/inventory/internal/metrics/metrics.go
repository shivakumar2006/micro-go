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
	InventoryRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "inventory_request_total",
			Help: "total inventory request total",
		},
		[]string{"method", "path", "status"},
	)

	InventoryRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "inventory_request_duration_seconds",
			Help: "inventory request duration per second",
		},
		[]string{"method", "path", "status"},
	)

	// internal vehicle service calls
	InventoryVehicleServiceCallTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "inventory_vehicle_service_call_total",
			Help: "inventory vehicle service call total",
		},
		[]string{"method", "path", "status"},
	)

	InventoryVehicleServiceCallSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "inventory_vehicle_service_call_success",
			Help: "inventory vehicle service call success",
		},
	)

	InventoryVehicleServiceCallFailure = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "inventory_vehicle_service_call_failure",
			Help: "inventory vehicle service call failure",
		},
	)

	InventoryVehicleServiceCallDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "inventory_vehicle_service_call_duration_seconds",
			Help: "inventory vehicle service call duration per second",
		},
		[]string{"method", "path", "status"},
	)

	// internal order service calls
	InventoryOrderServiceCallTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "inventory_order_service_call_total",
			Help: "inventory order service call total",
		},
	)

	InventoryOrderServiceCallSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "inventory_order_service_call_success",
			Help: "inventory order service call success",
		},
	)

	InventoryOrderServiceCallFailure = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "inventory_order_service_call_failure",
			Help: "inventory order service call failure",
		},
	)

	InventoryOrderServiceCallDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "inventory_order_service_call_duration_seconds",
			Help: "inventory order service call duration per second",
		},
		[]string{"method", "path", "status"},
	)

	// payment service events
	InventoryPaymentEventsProcessedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "inventory_payment_events_processed_total",
			Help: "Total number of payment events successfully processed",
		},
	)

	InventoryPaymentEventsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "inventory_payment_events_failed_total",
			Help: "Total number of payment events that failed to process",
		},
	)
)

func init() {
	prometheus.MustRegister(InventoryRequestTotal)
	prometheus.MustRegister(InventoryRequestDuration)
	prometheus.MustRegister(InventoryVehicleServiceCallTotal)
	prometheus.MustRegister(InventoryVehicleServiceCallDuration)
	prometheus.MustRegister(InventoryOrderServiceCallTotal)
	prometheus.MustRegister(InventoryOrderServiceCallDuration)
	prometheus.MustRegister(InventoryPaymentEventsProcessedTotal)
	prometheus.MustRegister(InventoryPaymentEventsFailedTotal)
}

func metricsMiddleware(next http.Handler) http.Handler {
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

		duration := time.Since(start).Seconds()

		status := strconv.Itoa(sw.StatusCode)

		InventoryRequestTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		InventoryRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
	})
}
