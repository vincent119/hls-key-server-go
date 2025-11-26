package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutConfig holds timeout middleware configuration
type TimeoutConfig struct {
	// Timeout duration for requests (default: 30 seconds)
	Timeout time.Duration
	// ErrorMessage returned when timeout occurs
	ErrorMessage string
}

// DefaultTimeoutConfig returns default timeout configuration
func DefaultTimeoutConfig() *TimeoutConfig {
	return &TimeoutConfig{
		Timeout:      30 * time.Second,
		ErrorMessage: "Request timeout",
	}
}

// Timeout returns a middleware that adds timeout control to requests
//
// Usage:
//
//	router.Use(middleware.Timeout(30 * time.Second))
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return TimeoutWithConfig(&TimeoutConfig{
		Timeout:      timeout,
		ErrorMessage: "Request timeout",
	})
}

// TimeoutWithConfig returns a middleware with custom timeout configuration
func TimeoutWithConfig(config *TimeoutConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultTimeoutConfig()
	}

	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}

	if config.ErrorMessage == "" {
		config.ErrorMessage = "Request timeout"
	}

	return func(c *gin.Context) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), config.Timeout)
		defer cancel()

		// Replace request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to signal completion
		finished := make(chan struct{}, 1)

		go func() {
			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-finished:
			// Request completed successfully
			return
		case <-ctx.Done():
			// Timeout occurred
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
				"error": config.ErrorMessage,
			})
			return
		}
	}
}
