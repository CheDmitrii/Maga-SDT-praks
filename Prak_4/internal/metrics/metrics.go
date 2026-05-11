package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HttpRequestsTotal counts all incoming HTTP requests by method and path.
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	// HttpErrorsTotal counts HTTP error responses (status >= 400) by method, path, and status code.
	HttpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_errors_total",
			Help: "Total number of HTTP error responses",
		},
		[]string{"method", "path", "status_code"},
	)

	// HttpRequestDuration tracks the distribution of HTTP request durations.
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// HttpRequestsTotal counts all incoming HTTP requests by students get method and path (additioanal task)
	HttpStudentRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_get_student_request_total",
			Help: "HTTP get student request",
		},
		[]string{"id"},
	)
)
