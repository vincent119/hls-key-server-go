package v1

import "github.com/gin-gonic/gin"

// GetRouteGroups is a function that returns all route groups
// @Summary Get all route groups
// @Description Get all route groups
// @Tags Route
func GetRouteGroups() []interface{ RegisterRoutes(*gin.RouterGroup) } {
	return []interface{ RegisterRoutes(*gin.RouterGroup) }{
		NewHlsKeyRoute(),
		NewAuthRoutes(),
	}
}
