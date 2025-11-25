package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"hls-key-server-go/internal/apperrors"
)

// mockHLSService implements a mock HLS service for testing
type mockHLSService struct {
	keys      map[string][]byte
	reloadErr error
}

func newMockHLSService() *mockHLSService {
	return &mockHLSService{
		keys: map[string][]byte{
			"test.key":   []byte("test-key-content-16b"),
			"stream.key": []byte("stream-key-data-16"),
		},
	}
}

func (m *mockHLSService) GetKey(_ context.Context, keyName string) ([]byte, error) {
	if keyName == "" {
		return nil, apperrors.ErrInvalidKeyName
	}
	key, ok := m.keys[keyName]
	if !ok {
		return nil, apperrors.ErrKeyNotFound
	}
	return key, nil
}

func (m *mockHLSService) ListKeys(_ context.Context) []string {
	keys := make([]string, 0, len(m.keys))
	for k := range m.keys {
		keys = append(keys, k)
	}
	return keys
}

func (m *mockHLSService) ReloadKeys(_ context.Context) error {
	return m.reloadErr
}

// HLSServiceInterface defines the interface for HLS service operations
type HLSServiceInterface interface {
	GetKey(ctx context.Context, keyName string) ([]byte, error)
	ListKeys(ctx context.Context) []string
	ReloadKeys(ctx context.Context) error
}

// testHLSHandler wraps HLSHandler for testing with mock service
type testHLSHandler struct {
	service HLSServiceInterface
	logger  *zap.Logger
}

func newTestHLSHandler(service HLSServiceInterface) *testHLSHandler {
	return &testHLSHandler{
		service: service,
		logger:  zap.NewNop(),
	}
}

func (h *testHLSHandler) GetKey(c *gin.Context) {
	keyName := c.Query("key")
	if keyName == "" {
		keyName = c.PostForm("key")
	}
	if keyName == "" {
		keyName = "stream.key"
	}

	keyData, err := h.service.GetKey(c.Request.Context(), keyName)
	if err != nil {
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

func (h *testHLSHandler) ListKeys(c *gin.Context) {
	keys := h.service.ListKeys(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

func (h *testHLSHandler) ReloadKeys(c *gin.Context) {
	if err := h.service.ReloadKeys(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload keys"})
		return
	}
	keys := h.service.ListKeys(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"message": "Keys reloaded successfully",
		"count":   len(keys),
	})
}

func setupTestRouter(handler *testHLSHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	hlsGroup := router.Group("/api/v1/hls")
	{
		hlsGroup.POST("/key", handler.GetKey)
		hlsGroup.GET("/keys", handler.ListKeys)
		hlsGroup.POST("/reload", handler.ReloadKeys)
	}
	return router
}

func TestHLSHandler_GetKey(t *testing.T) {
	mockService := newMockHLSService()
	handler := newTestHLSHandler(mockService)
	router := setupTestRouter(handler)

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		contentType    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "get key via POST form data",
			method:         "POST",
			url:            "/api/v1/hls/key",
			body:           "key=test.key",
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: http.StatusOK,
			expectedBody:   "test-key-content-16b",
		},
		{
			name:           "get key via query parameter",
			method:         "POST",
			url:            "/api/v1/hls/key?key=test.key",
			body:           "",
			contentType:    "",
			expectedStatus: http.StatusOK,
			expectedBody:   "test-key-content-16b",
		},
		{
			name:           "get default key when no key specified",
			method:         "POST",
			url:            "/api/v1/hls/key",
			body:           "",
			contentType:    "",
			expectedStatus: http.StatusOK,
			expectedBody:   "stream-key-data-16",
		},
		{
			name:           "key not found",
			method:         "POST",
			url:            "/api/v1/hls/key",
			body:           "key=nonexistent.key",
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"Key not found"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
				if tt.contentType != "" {
					req.Header.Set("Content-Type", tt.contentType)
				}
			} else {
				req = httptest.NewRequest(tt.method, tt.url, nil)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestHLSHandler_ListKeys(t *testing.T) {
	mockService := newMockHLSService()
	handler := newTestHLSHandler(mockService)
	router := setupTestRouter(handler)

	req := httptest.NewRequest("GET", "/api/v1/hls/keys", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "keys") {
		t.Errorf("expected response to contain 'keys', got %q", body)
	}
	if !strings.Contains(body, "test.key") {
		t.Errorf("expected response to contain 'test.key', got %q", body)
	}
}

func TestHLSHandler_ReloadKeys(t *testing.T) {
	tests := []struct {
		name           string
		reloadErr      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful reload",
			reloadErr:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Keys reloaded successfully",
		},
		{
			name:           "reload failure",
			reloadErr:      apperrors.ErrKeyReadFailed,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to reload keys",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := newMockHLSService()
			mockService.reloadErr = tt.reloadErr
			handler := newTestHLSHandler(mockService)
			router := setupTestRouter(handler)

			req := httptest.NewRequest("POST", "/api/v1/hls/reload", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestHLSHandler_GetKey_ContentType(t *testing.T) {
	mockService := newMockHLSService()
	handler := newTestHLSHandler(mockService)
	router := setupTestRouter(handler)

	form := url.Values{}
	form.Add("key", "test.key")
	req := httptest.NewRequest("POST", "/api/v1/hls/key", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/octet-stream" {
		t.Errorf("expected Content-Type 'application/octet-stream', got %q", contentType)
	}
}
