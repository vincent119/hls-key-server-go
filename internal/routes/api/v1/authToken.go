package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"hls-key-server-go/internal/configs"
	"hls-key-server-go/internal/handler/logging"
	"hls-key-server-go/internal/handler/middleware"
	"net/http"
)

type AuthRoutes struct{}

func NewAuthRoutes() *AuthRoutes {
	return &AuthRoutes{}
}

// AuthTokenHandler
// @Summary Generate auth token
// @Description Generates a JWT token if the username and custom header are valid
// @Tags Auth
// @Accept  application/x-www-form-urlencoded
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
		authGroup.POST("/token", func(c *gin.Context) {
			username := c.PostForm("username")

			// user check empty
			if username == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})

				logging.InitZapLogging().Info("User name is empty",
					zap.String("ip", c.ClientIP()),
				)
				return
			}

			// user check
			if username != configs.Conf.JwtSecret.User {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username"})
				logging.InitZapLogging().Info("User name is Invalid",
					zap.String("ip", c.ClientIP()),
				)
				return
			}

			// check custom http header
			headerKey := configs.Conf.JwtSecret.HeaderKey
			headerValue := configs.Conf.JwtSecret.HeaderValue

			receivedHeader := c.GetHeader(headerKey)
			if receivedHeader == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Custom header is required"})

				logging.InitZapLogging().Info("Custom header is empty",
					zap.String("ip", c.ClientIP()),
				)

				return
			}

			if receivedHeader != headerValue {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Custom header is invalid"})

				logging.InitZapLogging().Info("Custom header value mot equal",
					zap.String("ip", c.ClientIP()),
				)
				return
			}

			userIP := c.ClientIP()
			c.Set("user_ip", userIP) // set user ip to context
			fmt.Println("user ip: ", userIP)
			// generate token
			token, err := middleware.GenerateJWT(username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})

				logging.InitZapLogging().Info("Failed to generate token",
					zap.String("ip", c.ClientIP()),
				)
				return
			}
			// return token
			c.JSON(http.StatusOK, gin.H{"token": token})
		})
	}
}
