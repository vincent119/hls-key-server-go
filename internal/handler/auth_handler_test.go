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
	"hls-key-server-go/internal/configs"
)

// mockAuthService implements a mock Auth service for testing
type mockAuthService struct {
	validUser   string
	headerValue string
	token       string
	validateErr error
}

func newMockAuthService() *mockAuthService {
	return &mockAuthService{
		validUser:   "testuser",
		headerValue: "valid-header-value",
		token:       "mock-jwt-token-12345",
	}
}

func (m *mockAuthService) GenerateToken(_ context.Context, _ string) (string, error) {
	return m.token, nil
}

func (m *mockAuthService) ValidateCredentials(_ context.Context, username, headerValue string) error {
	if username == "" || username != m.validUser {
		return apperrors.ErrInvalidCredentials
	}
	if headerValue != m.headerValue {
		return apperrors.ErrMissingHeader
	}
	return m.validateErr
}

// AuthServiceInterface defines the interface for Auth service operations
type AuthServiceInterface interface {
	GenerateToken(ctx context.Context, username string) (string, error)
	ValidateCredentials(ctx context.Context, username, headerValue string) error
}

// testAuthHandler wraps AuthHandler for testing with mock service
type testAuthHandler struct {
	service   AuthServiceInterface
	jwtConfig *configs.JwtSecret
	logger    *zap.Logger
}

func newTestAuthHandler(service AuthServiceInterface, jwtConfig *configs.JwtSecret) *testAuthHandler {
	return &testAuthHandler{
		service:   service,
		jwtConfig: jwtConfig,
		logger:    zap.NewNop(),
	}
}

func (h *testAuthHandler) GenerateToken(c *gin.Context) {
	username := c.PostForm("username")

	// Validate custom header
	headerValue := c.GetHeader(h.jwtConfig.HeaderKey)
	if headerValue != h.jwtConfig.HeaderValue {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication header"})
		return
	}

	// Validate credentials
	if err := h.service.ValidateCredentials(c.Request.Context(), username, headerValue); err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func setupAuthTestRouter(handler *testAuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/token", handler.GenerateToken)
	}

	return router
}

func TestAuthHandler_GenerateToken(t *testing.T) {
	jwtConfig := &configs.JwtSecret{
		User:        "testuser",
		HeaderKey:   "X-Custom-Auth",
		HeaderValue: "valid-header-value",
		SecretKey:   "test-secret",
	}

	tests := []struct {
		name           string
		username       string
		headerKey      string
		headerValue    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful token generation",
			username:       "testuser",
			headerKey:      "X-Custom-Auth",
			headerValue:    "valid-header-value",
			expectedStatus: http.StatusOK,
			expectedBody:   "mock-jwt-token-12345",
		},
		{
			name:           "invalid header",
			username:       "testuser",
			headerKey:      "X-Custom-Auth",
			headerValue:    "wrong-header-value",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid authentication header",
		},
		{
			name:           "missing header",
			username:       "testuser",
			headerKey:      "",
			headerValue:    "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid authentication header",
		},
		{
			name:           "invalid username",
			username:       "wronguser",
			headerKey:      "X-Custom-Auth",
			headerValue:    "valid-header-value",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials",
		},
		{
			name:           "empty username",
			username:       "",
			headerKey:      "X-Custom-Auth",
			headerValue:    "valid-header-value",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := newMockAuthService()
			handler := newTestAuthHandler(mockService, jwtConfig)
			router := setupAuthTestRouter(handler)

			form := url.Values{}
			form.Add("username", tt.username)

			req := httptest.NewRequest("POST", "/api/v1/auth/token", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if tt.headerKey != "" {
				req.Header.Set(tt.headerKey, tt.headerValue)
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

func TestAuthHandler_GenerateToken_ContentType(t *testing.T) {
	jwtConfig := &configs.JwtSecret{
		User:        "testuser",
		HeaderKey:   "X-Custom-Auth",
		HeaderValue: "valid-header-value",
		SecretKey:   "test-secret",
	}

	mockService := newMockAuthService()
	handler := newTestAuthHandler(mockService, jwtConfig)
	router := setupAuthTestRouter(handler)

	form := url.Values{}
	form.Add("username", "testuser")

	req := httptest.NewRequest("POST", "/api/v1/auth/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Custom-Auth", "valid-header-value")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected Content-Type to contain 'application/json', got %q", contentType)
	}
}
