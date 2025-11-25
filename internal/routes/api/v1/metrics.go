// Package v1 provides HTTP route registration for API version 1 endpoints.
package v1

import (
	"github.com/gin-gonic/gin"

	"hls-key-server-go/internal/handler"
)

// MetricsRoute handles metrics route registration
type MetricsRoute struct {
	metricsHandler *handler.MetricsHandler
}

// NewMetricsRoute creates a new metrics route
func NewMetricsRoute(metricsHandler *handler.MetricsHandler) *MetricsRoute {
	return &MetricsRoute{
		metricsHandler: metricsHandler,
	}
}

// RegisterRoutes registers metrics routes
func (r *MetricsRoute) RegisterRoutes(rg *gin.RouterGroup) {
	// Register at parent router level (not under /api/v1)
	router := rg.Group("")
	router.GET("/metrics", r.metricsHandler.BasicAuth(), r.metricsHandler.Handler())
}
