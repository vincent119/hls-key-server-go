package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestTimeout(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		timeout        time.Duration
		handlerDelay   time.Duration
		expectedStatus int
		shouldTimeout  bool
	}{
		{
			name:           "request completes before timeout",
			timeout:        100 * time.Millisecond,
			handlerDelay:   10 * time.Millisecond,
			expectedStatus: http.StatusOK,
			shouldTimeout:  false,
		},
		{
			name:           "request times out",
			timeout:        50 * time.Millisecond,
			handlerDelay:   200 * time.Millisecond,
			expectedStatus: http.StatusGatewayTimeout,
			shouldTimeout:  true,
		},
		{
			name:           "instant response",
			timeout:        1 * time.Second,
			handlerDelay:   0,
			expectedStatus: http.StatusOK,
			shouldTimeout:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			router := gin.New()
			router.Use(Timeout(tt.timeout))
			router.GET("/test", func(c *gin.Context) {
				if tt.handlerDelay > 0 {
					select {
					case <-time.After(tt.handlerDelay):
						c.JSON(http.StatusOK, gin.H{"status": "ok"})
					case <-c.Request.Context().Done():
						return
					}
				} else {
					c.JSON(http.StatusOK, gin.H{"status": "ok"})
				}
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestTimeoutWithConfig(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	config := &TimeoutConfig{
		Timeout:      100 * time.Millisecond,
		ErrorMessage: "Custom timeout message",
	}

	router := gin.New()
	router.Use(TimeoutWithConfig(config))
	router.GET("/test", func(c *gin.Context) {
		time.Sleep(200 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("expected status %d, got %d", http.StatusGatewayTimeout, w.Code)
	}
}

func TestDefaultTimeoutConfig(t *testing.T) {
	config := DefaultTimeoutConfig()

	if config.Timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", config.Timeout)
	}

	if config.ErrorMessage == "" {
		t.Error("expected non-empty error message")
	}
}

func TestTimeoutContextPropagation(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(Timeout(100 * time.Millisecond))
	router.GET("/test", func(c *gin.Context) {
		ctx := c.Request.Context()
		if ctx == nil {
			t.Error("expected context to be set")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no context"})
			return
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Error("expected context to have deadline")
		}

		if time.Until(deadline) > 100*time.Millisecond {
			t.Errorf("expected deadline within 100ms, got %v", time.Until(deadline))
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func BenchmarkTimeout(b *testing.B) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(Timeout(1 * time.Second))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}
	})
}
