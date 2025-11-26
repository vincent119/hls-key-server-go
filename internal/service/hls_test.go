package service

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"

	"hls-key-server-go/internal/apperrors"
)

// mockKeyRepository implements repository.KeyRepository for testing
type mockKeyRepository struct {
	keys      map[string][]byte
	getErr    error
	reloadErr error
	listKeys  []string
}

func newMockKeyRepository() *mockKeyRepository {
	return &mockKeyRepository{
		keys: map[string][]byte{
			"test.key":   []byte("test-key-content-1234567890123456"),
			"stream.key": []byte("stream-key-data-1234567890123456"),
		},
	}
}

func (m *mockKeyRepository) Get(_ context.Context, name string) ([]byte, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	data, ok := m.keys[name]
	if !ok {
		return nil, apperrors.ErrKeyNotFound
	}
	// Return a copy to simulate real repository behavior
	return append([]byte(nil), data...), nil
}

func (m *mockKeyRepository) List(_ context.Context) []string {
	if m.listKeys != nil {
		return m.listKeys
	}
	keys := make([]string, 0, len(m.keys))
	for k := range m.keys {
		keys = append(keys, k)
	}
	return keys
}

func (m *mockKeyRepository) Reload(_ context.Context) error {
	return m.reloadErr
}

func TestNewHLSService(t *testing.T) {
	repo := newMockKeyRepository()
	logger := zap.NewNop()
	service := NewHLSService(repo, logger)

	if service == nil {
		t.Fatal("NewHLSService() returned nil")
	}
}

func TestHLSService_GetKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		keyName     string
		wantErr     bool
		wantErrType error
		setupRepo   func() *mockKeyRepository
	}{
		{
			name:      "get existing key",
			keyName:   "test.key",
			wantErr:   false,
			setupRepo: newMockKeyRepository,
		},
		{
			name:      "get another existing key",
			keyName:   "stream.key",
			wantErr:   false,
			setupRepo: newMockKeyRepository,
		},
		{
			name:        "get non-existent key",
			keyName:     "missing.key",
			wantErr:     true,
			wantErrType: apperrors.ErrKeyNotFound,
			setupRepo:   newMockKeyRepository,
		},
		{
			name:        "repository error",
			keyName:     "test.key",
			wantErr:     true,
			wantErrType: apperrors.ErrInvalidKeyName,
			setupRepo: func() *mockKeyRepository {
				repo := newMockKeyRepository()
				repo.getErr = apperrors.ErrInvalidKeyName
				return repo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.setupRepo()
			logger := zap.NewNop()
			service := NewHLSService(repo, logger)
			ctx := context.Background()

			data, err := service.GetKey(ctx, tt.keyName)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetKey() error = nil, want error")
					return
				}
				// Verify error wrapping
				if tt.wantErrType != nil && !errors.Is(err, tt.wantErrType) {
					t.Errorf("GetKey() error type = %T, want %T", err, tt.wantErrType)
				}
				return
			}

			if err != nil {
				t.Errorf("GetKey() unexpected error = %v", err)
				return
			}
			if len(data) == 0 {
				t.Error("GetKey() returned empty data")
			}
			// Verify key content length matches expected AES key size
			if len(data) < 16 {
				t.Errorf("GetKey() returned key with length %d, want at least 16 bytes", len(data))
			}
		})
	}
}

func TestHLSService_ListKeys(t *testing.T) {
	repo := newMockKeyRepository()
	logger := zap.NewNop()
	service := NewHLSService(repo, logger)

	ctx := context.Background()
	keys := service.ListKeys(ctx)

	if len(keys) == 0 {
		t.Error("ListKeys() returned empty list")
	}
}

func TestHLSService_ReloadKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		wantErr   bool
		setupRepo func() *mockKeyRepository
	}{
		{
			name:    "successful reload",
			wantErr: false,
			setupRepo: func() *mockKeyRepository {
				return newMockKeyRepository()
			},
		},
		{
			name:    "reload with error",
			wantErr: true,
			setupRepo: func() *mockKeyRepository {
				repo := newMockKeyRepository()
				repo.reloadErr = apperrors.ErrKeyNotFound
				return repo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.setupRepo()
			logger := zap.NewNop()
			service := NewHLSService(repo, logger)
			ctx := context.Background()

			err := service.ReloadKeys(ctx)

			if tt.wantErr {
				if err == nil {
					t.Error("ReloadKeys() error = nil, want error")
				}
				return
			}

			if err != nil {
				t.Errorf("ReloadKeys() unexpected error = %v", err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkHLSService_GetKey(b *testing.B) {
	repo := newMockKeyRepository()
	logger := zap.NewNop()
	service := NewHLSService(repo, logger)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = service.GetKey(ctx, "test.key")
		}
	})
}

func BenchmarkHLSService_ListKeys(b *testing.B) {
	repo := newMockKeyRepository()
	logger := zap.NewNop()
	service := NewHLSService(repo, logger)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.ListKeys(ctx)
	}
}
