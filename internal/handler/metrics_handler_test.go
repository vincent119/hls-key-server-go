// Package handler provides HTTP request handlers including metrics endpoint.
package handler

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"hls-key-server-go/internal/configs"
)

func TestMetricsHandler_BasicAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		user           string
		password       string
		expectStatus   int
		expectAuthCall bool
	}{
		{
			name:           "valid credentials",
			user:           "admin",
			password:       "sshhsuuwgwhysgs",
			expectStatus:   http.StatusOK,
			expectAuthCall: false,
		},
		{
			name:           "invalid user",
			user:           "wrong",
			password:       "sshhsuuwgwhysgs",
			expectStatus:   http.StatusUnauthorized,
			expectAuthCall: true,
		},
		{
			name:           "invalid password",
			user:           "admin",
			password:       "wrong",
			expectStatus:   http.StatusUnauthorized,
			expectAuthCall: true,
		},
		{
			name:           "no credentials",
			user:           "",
			password:       "",
			expectStatus:   http.StatusUnauthorized,
			expectAuthCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			cfg := &configs.Config{
				Metric: configs.Metric{
					User:     "admin",
					Password: "sshhsuuwgwhysgs",
				},
			}

			logger := zap.NewNop()
			handler := NewMetricsHandler(cfg, logger)

			router := gin.New()
			router.GET("/metrics", handler.BasicAuth(), func(c *gin.Context) {
				c.String(http.StatusOK, "metrics")
			})

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
			if tt.user != "" || tt.password != "" {
				auth := base64.StdEncoding.EncodeToString([]byte(tt.user + ":" + tt.password))
				req.Header.Set("Authorization", "Basic "+auth)
			}

			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			if w.Code != tt.expectStatus {
				t.Errorf("expected status %d, got %d", tt.expectStatus, w.Code)
			}

			if tt.expectAuthCall && w.Header().Get("WWW-Authenticate") == "" {
				t.Error("expected WWW-Authenticate header to be set")
			}
		})
	}
}

func TestMetricsHandler_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &configs.Config{
		Metric: configs.Metric{
			User:     "admin",
			Password: "testpass",
		},
	}

	logger := zap.NewNop()
	handler := NewMetricsHandler(cfg, logger)

	router := gin.New()
	router.GET("/metrics", handler.BasicAuth(), handler.Handler())

	// Create request with valid auth
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	auth := base64.StdEncoding.EncodeToString([]byte("admin:testpass"))
	req.Header.Set("Authorization", "Basic "+auth)

	w := httptest.NewRecorder()

	// Execute
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Check that response contains prometheus metrics
	body := w.Body.String()
	if body == "" {
		t.Error("expected non-empty metrics response")
	}
}
