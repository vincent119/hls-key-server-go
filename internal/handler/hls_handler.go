// Package handler provides HTTP request handlers using service layer
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"hls-key-server-go/internal/apperrors"
	"hls-key-server-go/internal/service"
)

// HLSHandler handles HLS key requests
type HLSHandler struct {
	service *service.HLSService
	logger  *zap.Logger
}

// NewHLSHandler creates a new HLS handler
func NewHLSHandler(service *service.HLSService, logger *zap.Logger) *HLSHandler {
	return &HLSHandler{
		service: service,
		logger:  logger,
	}
}

// GetKey handles the key retrieval request
// @Summary Get encryption key
// @Description Retrieves an HLS encryption key by name
// @Tags HLS
// @Accept json
// @Produce octet-stream
// @Param key query string false "Key name (default: stream.key)"
// @Security BearerAuth
// @Success 200 {file} binary "Encryption key"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Key not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/v1/hls/key [post]
func (h *HLSHandler) GetKey(c *gin.Context) {
	// Support both query parameter and form data
	keyName := c.Query("key")
	if keyName == "" {
		keyName = c.PostForm("key")
	}
	if keyName == "" {
		keyName = "stream.key"
	}

	h.logger.Info("key request",
		zap.String("key", keyName),
		zap.String("ip", c.ClientIP()),
	)

	keyData, err := h.service.GetKey(c.Request.Context(), keyName)
	if err != nil {
		h.logger.Error("failed to get key",
			zap.String("key", keyName),
			zap.Error(err),
		)

		if apperrors.IsKeyNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
			return
		}
		if apperrors.IsInvalidKeyName(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key name"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve key"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", keyData)
}

// ListKeys handles listing all available keys
// @Summary List all keys
// @Description Lists all available encryption keys
// @Tags HLS
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string][]string "List of keys"
// @Router /api/v1/hls/keys [get]
func (h *HLSHandler) ListKeys(c *gin.Context) {
	keys := h.service.ListKeys(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// ReloadKeys handles reloading all keys from storage
// @Summary Reload keys
// @Description Reloads all encryption keys from the filesystem
// @Tags HLS
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Reload status"
// @Failure 500 {object} map[string]string "Server error"
// @Router /api/v1/hls/reload [post]
func (h *HLSHandler) ReloadKeys(c *gin.Context) {
	h.logger.Info("key reload request",
		zap.String("ip", c.ClientIP()),
	)

	if err := h.service.ReloadKeys(c.Request.Context()); err != nil {
		h.logger.Error("failed to reload keys", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload keys"})
		return
	}

	keys := h.service.ListKeys(c.Request.Context())
	h.logger.Info("keys reloaded successfully", zap.Int("count", len(keys)))
	c.JSON(http.StatusOK, gin.H{
		"message": "Keys reloaded successfully",
		"count":   len(keys),
	})
}
