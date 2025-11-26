package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestHealthCheck(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkResponse  bool
	}{
		{
			name:           "GET health check",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "POST health check should also work",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			checkResponse:  true,
		},
		{
			name:           "HEAD health check",
			method:         http.MethodHead,
			expectedStatus: http.StatusOK,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := gin.New()
			router.Any("/healthz", HealthCheck)

			req := httptest.NewRequest(tt.method, "/healthz", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse {
				var response ResponseHealthCheck
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if response.Status != "OK" {
					t.Errorf("expected status 'OK', got %q", response.Status)
				}

				if response.RecvTime == "" {
					t.Error("expected RecvTime to be set")
				}

				if response.RecvTimeUTC == "" {
					t.Error("expected RecvTimeUTC to be set")
				}

				// Verify time format
				if _, err := time.Parse("2006-01-02T15:04:05", response.RecvTime); err != nil {
					t.Errorf("invalid RecvTime format: %v", err)
				}

				if _, err := time.Parse(time.RFC3339, response.RecvTimeUTC); err != nil {
					t.Errorf("invalid RecvTimeUTC format: %v", err)
				}
			}
		})
	}
}

func TestHealthCheck_ContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/healthz", HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", contentType)
	}
}

func BenchmarkHealthCheck(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/healthz", HealthCheck)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}
