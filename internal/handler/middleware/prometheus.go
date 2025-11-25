// Package middleware provides HTTP middleware functions
package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"hls-key-server-go/internal/pkg/metrics"
)

// PrometheusMiddleware records HTTP request metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Track concurrent connections
		metrics.ConcurrentConnections.Inc()
		defer metrics.ConcurrentConnections.Dec()

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// Record request count
		metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()

		// Record request duration
		metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
