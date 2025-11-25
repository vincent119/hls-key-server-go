// Package metrics provides custom Prometheus metrics for the HLS key server
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal tracks total HTTP requests by method, path and status
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hls_http_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestDuration tracks HTTP request duration
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "hls_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// KeyRequestsTotal tracks total key requests
	KeyRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hls_key_requests_total",
			Help: "Total number of HLS key requests",
		},
		[]string{"key_name", "status"},
	)

	// KeyCacheHits tracks cache hit rate
	KeyCacheHits = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "hls_key_cache_hits_total",
			Help: "Total number of key cache hits",
		},
	)

	// KeyCacheMisses tracks cache misses
	KeyCacheMisses = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "hls_key_cache_misses_total",
			Help: "Total number of key cache misses",
		},
	)

	// ActiveKeys tracks the number of active keys
	ActiveKeys = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "hls_active_keys",
			Help: "Number of currently active HLS keys",
		},
	)

	// AuthAttempts tracks authentication attempts
	AuthAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hls_auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"result"},
	)

	// TokenGenerations tracks JWT token generations
	TokenGenerations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "hls_token_generations_total",
			Help: "Total number of JWT tokens generated",
		},
	)

	// TokenValidations tracks JWT token validations
	TokenValidations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hls_token_validations_total",
			Help: "Total number of JWT token validations",
		},
		[]string{"result"},
	)

	// ErrorsTotal tracks errors by type
	ErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hls_errors_total",
			Help: "Total number of errors by type",
		},
		[]string{"type"},
	)

	// KeyReloadDuration tracks key reload duration
	KeyReloadDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "hls_key_reload_duration_seconds",
			Help:    "Duration of key reload operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	// ConcurrentConnections tracks current concurrent connections
	ConcurrentConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "hls_concurrent_connections",
			Help: "Number of concurrent HTTP connections",
		},
	)

	// ServerUptime tracks server uptime
	ServerUptime = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "hls_server_uptime_seconds",
			Help: "Server uptime in seconds",
		},
	)

	// KeyFileSize tracks the size of key files
	KeyFileSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hls_key_file_size_bytes",
			Help: "Size of HLS key files in bytes",
		},
		[]string{"key_name"},
	)

	// APIVersion tracks API version info
	APIVersion = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hls_api_version_info",
			Help: "API version information",
		},
		[]string{"version", "mode"},
	)
)

// Init initializes custom metrics with default values
func Init(version, mode string) {
	// Set API version info
	APIVersion.WithLabelValues(version, mode).Set(1)
}
