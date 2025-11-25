// Package v1 provides HTTP route registration for API version 1 endpoints.
package v1

import (
	"hls-key-server-go/internal/handler"

	"github.com/gin-gonic/gin"
)

// GetRouteGroups is a function that returns all route groups
// @Summary Get all route groups
// @Description Get all route groups
// @Tags Route
func GetRouteGroups(hlsHandler *handler.HLSHandler, authHandler *handler.AuthHandler, metricsHandler *handler.MetricsHandler) []interface{ RegisterRoutes(*gin.RouterGroup) } {
	return []interface{ RegisterRoutes(*gin.RouterGroup) }{
		NewHlsKeyRoute(hlsHandler),
		NewAuthRoutes(authHandler),
		NewMetricsRoute(metricsHandler),
	}
}
