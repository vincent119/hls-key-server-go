

// Package handler provides HTTP request handlers including metrics endpoint.
package handler

import (
	"crypto/subtle"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"hls-key-server-go/internal/configs"
)

// MetricsHandler handles Prometheus metrics endpoint with basic auth
type MetricsHandler struct {
	config *configs.Config
	logger *zap.Logger
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(config *configs.Config, logger *zap.Logger) *MetricsHandler {
	return &MetricsHandler{
		config: config,
		logger: logger,
	}
}

// BasicAuth provides basic authentication middleware for metrics endpoint
func (h *MetricsHandler) BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok {
			h.logger.Warn("metrics endpoint accessed without basic auth")
			c.Header("WWW-Authenticate", `Basic realm="metrics"`)
			c.AbortWithStatus(401)
			return
		}

		// Use constant time comparison to prevent timing attacks
		userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(h.config.Metric.User)) == 1
		passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(h.config.Metric.Password)) == 1

		if !userMatch || !passMatch {
			h.logger.Warn("metrics endpoint accessed with invalid credentials",
				zap.String("user", user),
			)
			c.Header("WWW-Authenticate", `Basic realm="metrics"`)
			c.AbortWithStatus(401)
			return
		}

		c.Next()
	}
}

// Handler returns the Prometheus metrics handler
// @Summary Prometheus metrics endpoint
// @Description Returns Prometheus metrics with basic authentication
// @Tags metrics
// @Security BasicAuth
// @Success 200 {string} string "Prometheus metrics in text format"
// @Failure 401 {string} string "Unauthorized"
// @Router /metrics [get]
func (h *MetricsHandler) Handler() gin.HandlerFunc {
	handler := promhttp.Handler()
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
