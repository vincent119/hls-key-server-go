package service

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"hls-key-server-go/internal/configs"
)

func TestNewAuthService(t *testing.T) {
	config := &configs.JwtSecret{
		SecretKey: "test-secret-key",
		Expire:    10,
		User:      "testuser",
		Iss:       "test-issuer",
		Aud:       "test-audience",
	}
	logger := zap.NewNop()

	service := NewAuthService(config, logger)

	if service == nil {
		t.Fatal("NewAuthService() returned nil")
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	config := &configs.JwtSecret{
		SecretKey: "test-secret-key-for-jwt",
		Expire:    10,
		User:      "testuser",
		Iss:       "test-issuer",
		Aud:       "test-audience",
	}
	logger := zap.NewNop()

	service := NewAuthService(config, logger)

	ctx := context.Background()
	token, err := service.GenerateToken(ctx, "testuser")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}
}

func TestAuthService_ValidateCredentials(t *testing.T) {
	config := &configs.JwtSecret{
		SecretKey:   "test-secret",
		User:        "validuser",
		HeaderValue: "valid-header-value",
	}
	logger := zap.NewNop()

	service := NewAuthService(config, logger)

	tests := []struct {
		name        string
		username    string
		headerValue string
		wantErr     bool
	}{
		{
			name:        "valid credentials",
			username:    "validuser",
			headerValue: "valid-header-value",
			wantErr:     false,
		},
		{
			name:        "invalid username",
			username:    "invaliduser",
			headerValue: "valid-header-value",
			wantErr:     true,
		},
		{
			name:        "empty username",
			username:    "",
			headerValue: "valid-header-value",
			wantErr:     true,
		},
		{
			name:        "invalid header",
			username:    "validuser",
			headerValue: "invalid-header",
			wantErr:     true,
		},
		{
			name:        "empty header",
			username:    "validuser",
			headerValue: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := service.ValidateCredentials(ctx, tt.username, tt.headerValue)

			if tt.wantErr {
				if err == nil {
					t.Error("ValidateCredentials() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateCredentials() unexpected error = %v", err)
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	config := &configs.JwtSecret{
		SecretKey: "test-secret-key-for-validation",
		Expire:    10,
		User:      "testuser",
		Iss:       "test-issuer",
		Aud:       "test-audience",
	}
	logger := zap.NewNop()

	service := NewAuthService(config, logger)

	// Generate a valid token
	ctx := context.Background()
	token, err := service.GenerateToken(ctx, "testuser")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   token,
			wantErr: false,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "invalid.jwt.token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			claims, err := service.ValidateToken(ctx, tt.token)

			if tt.wantErr {
				if err == nil {
					t.Error("ValidateToken() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateToken() unexpected error = %v", err)
				return
			}

			if claims == nil {
				t.Error("ValidateToken() returned nil claims")
			}

			// Verify claims
			if sub, ok := claims["sub"].(string); !ok || sub != "testuser" {
				t.Errorf("ValidateToken() sub claim = %v, want testuser", sub)
			}
		})
	}
}

func BenchmarkAuthService_GenerateToken(b *testing.B) {
	config := &configs.JwtSecret{
		SecretKey: "benchmark-secret-key",
		Expire:    10,
		Iss:       "benchmark-issuer",
		Aud:       "benchmark-audience",
	}
	logger := zap.NewNop()
	service := NewAuthService(config, logger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GenerateToken(ctx, "testuser") //nolint:errcheck // benchmark
	}
}

func BenchmarkAuthService_ValidateToken(b *testing.B) {
	config := &configs.JwtSecret{
		SecretKey: "benchmark-secret-key",
		Expire:    10,
		Iss:       "benchmark-issuer",
		Aud:       "benchmark-audience",
	}
	logger := zap.NewNop()
	service := NewAuthService(config, logger)
	ctx := context.Background()

	token, err := service.GenerateToken(ctx, "testuser")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateToken(ctx, token) //nolint:errcheck // benchmark
	}
}
