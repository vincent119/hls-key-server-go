// Package handler provides HTTP request handlers using service layer
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"hls-key-server-go/internal/apperrors"
	"hls-key-server-go/internal/configs"
	"hls-key-server-go/internal/service"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	service   *service.AuthService
	jwtConfig *configs.JwtSecret
	logger    *zap.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service *service.AuthService, jwtConfig *configs.JwtSecret, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		service:   service,
		jwtConfig: jwtConfig,
		logger:    logger,
	}
}

// GenerateToken handles JWT token generation
// @Summary Generate auth token
// @Description Generates a JWT token if the username and custom header are valid
// @Tags Auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param header-key header string true "Custom authentication header"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/v1/auth/token [post]
func (h *AuthHandler) GenerateToken(c *gin.Context) {
	username := c.PostForm("username")

	h.logger.Info("token generation request",
		zap.String("username", username),
		zap.String("ip", c.ClientIP()),
	)

	// Validate custom header
	headerValue := c.GetHeader(h.jwtConfig.HeaderKey)
	if headerValue != h.jwtConfig.HeaderValue {
		h.logger.Warn("invalid custom header",
			zap.String("username", username),
			zap.String("ip", c.ClientIP()),
			zap.String("expected_header", h.jwtConfig.HeaderKey),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication header"})
		return
	}

	// Validate credentials
	if err := h.service.ValidateCredentials(c.Request.Context(), username, h.jwtConfig.HeaderValue); err != nil{
		h.logger.Warn("invalid credentials",
			zap.String("username", username),
			zap.Error(err),
		)

		if apperrors.IsInvalidCredentials(err) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	// Generate token
	token, err := h.service.GenerateToken(c.Request.Context(), username)
	if err != nil {
		h.logger.Error("failed to generate token",
			zap.String("username", username),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	h.logger.Info("token generated successfully",
		zap.String("username", username),
	)

	c.JSON(http.StatusOK, gin.H{"token": token})
}
