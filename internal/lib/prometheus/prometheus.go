package prometheus

import (
	prometheuslib "github.com/prometheus/client_golang/prometheus"
)

var (
	TotalRequests = prometheuslib.NewCounterVec(
		prometheuslib.CounterOpts{
			Name: "requests_total",
			Help: "Total number of HTTP requests processed by the API Gateway.",
		},
		[]string{"method", "route", "status"},
	)

	RequestDuration = prometheuslib.NewHistogramVec(
		prometheuslib.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Histogram of HTTP request durations in seconds.",
			Buckets: prometheuslib.DefBuckets, // Default buckets (e.g., 0.005s, 0.01s, 0.025s, etc.)
		},
		[]string{"method", "route"},
	)

	ErrorCount = prometheuslib.NewCounterVec(
		prometheuslib.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors encountered by the API Gateway.",
		},
		[]string{"type", "route"},
	)

	ServiceHealth = prometheuslib.NewGauge(prometheuslib.GaugeOpts{
		Name: "service_health_status",
		Help: "Health status of the service: 1=Healthy, 0=Unhealthy",
	})
)

func RegisterMetrics() {
	prometheuslib.MustRegister(
		TotalRequests,
		ErrorCount,
		RequestDuration,
		ServiceHealth,
	)
}
