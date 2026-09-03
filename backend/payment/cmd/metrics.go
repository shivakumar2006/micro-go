package main

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var paymentRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "payment_requests_total",
		Help: "total number of payment requests",
	},
	[]string{"method", "path", "status"},
)

type StatusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (s *StatusResponseWriter) WriteHeader(code int) {
	s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *StatusResponseWriter) Write(b []byte) (int, error) {
	if s.statusCode == 0 {
		s.statusCode = http.StatusOK
	}

	return s.ResponseWriter.Write(b)
}

func init() {
	prometheus.MustRegister(paymentRequests)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &StatusResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(sw, r)

		paymentRequests.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(sw.statusCode)).Inc()
	})
}
