package v1

import (
	"github.com/gin-gonic/gin"

	"hls-key-server-go/internal/handler"
)

// AuthRoutes handles authentication routes
type AuthRoutes struct {
	authHandler *handler.AuthHandler
}

// NewAuthRoutes creates a new auth routes handler
func NewAuthRoutes(authHandler *handler.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		authHandler: authHandler,
	}
}

// RegisterRoutes registers authentication routes
// @Summary Generate auth token
// @Description Generates a JWT token if the username and custom header are valid
// @Tags Auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param header-key header string true "Custom authentication header"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/v1/auth/token [post]
func (a *AuthRoutes) RegisterRoutes(group *gin.RouterGroup) {
	authGroup := group.Group("/auth")
	{
		authGroup.POST("/token", a.authHandler.GenerateToken)
	}
}
