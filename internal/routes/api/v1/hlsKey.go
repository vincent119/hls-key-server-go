package v1

import (
	"github.com/gin-gonic/gin"
	"hls-key-server-go/internal/handler/hls"
	"hls-key-server-go/internal/handler/middleware"
	//"net/http"
)

type HlsKeyRoute struct{}

func NewHlsKeyRoute() *HlsKeyRoute {
	return &HlsKeyRoute{}
}

// RegisterRoutes is a function that register routes for hls key
// @Summary Hls key
// @Description Hls key
// @Tags Hls
// @Accept  json
func (a *HlsKeyRoute) RegisterRoutes(group *gin.RouterGroup) {
	hlsGroup := group.Group("/hls")
	hlsGroup.Use(middleware.AuthMiddleware())
	{
		hlsGroup.POST("/key", hls.KeyHandler)
	}
}
