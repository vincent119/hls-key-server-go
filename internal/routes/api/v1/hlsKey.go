package v1

import (
	"github.com/gin-gonic/gin"

	"hls-key-server-go/internal/handler"
)

// HlsKeyRoute handles HLS key routes
type HlsKeyRoute struct {
	hlsHandler *handler.HLSHandler
}

// NewHlsKeyRoute creates a new HLS key route
func NewHlsKeyRoute(hlsHandler *handler.HLSHandler) *HlsKeyRoute {
	return &HlsKeyRoute{
		hlsHandler: hlsHandler,
	}
}

// RegisterRoutes registers HLS routes
// @Summary Hls key
// @Description Hls key
// @Tags Hls
// @Accept  json
func (a *HlsKeyRoute) RegisterRoutes(group *gin.RouterGroup) {
	hlsGroup := group.Group("/hls")
	{
		hlsGroup.POST("/key", a.hlsHandler.GetKey)
		hlsGroup.GET("/keys", a.hlsHandler.ListKeys)
		hlsGroup.POST("/reload", a.hlsHandler.ReloadKeys)
	}
}
