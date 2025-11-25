package service

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"hls-key-server-go/internal/apperrors"
)

// mockKeyRepository implements repository.KeyRepository for testing
type mockKeyRepository struct {
	keys   map[string][]byte
	getErr error
}

func newMockKeyRepository() *mockKeyRepository {
	return &mockKeyRepository{
		keys: map[string][]byte{
			"test.key": []byte("test-key-content"),
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
	return data, nil
}

func (m *mockKeyRepository) List(_ context.Context) []string {
	keys := make([]string, 0, len(m.keys))
	for k := range m.keys {
		keys = append(keys, k)
	}
	return keys
}

func (m *mockKeyRepository) Reload(_ context.Context) error {
	return nil
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
	repo := newMockKeyRepository()
	logger := zap.NewNop()
	service := NewHLSService(repo, logger)

	tests := []struct {
		name    string
		keyName string
		wantErr error
		setup   func()
	}{
		{
			name:    "get existing key",
			keyName: "test.key",
			wantErr: nil,
			setup:   func() {},
		},
		{
			name:    "get non-existent key",
			keyName: "missing.key",
			wantErr: apperrors.ErrKeyNotFound,
			setup:   func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			ctx := context.Background()
			data, err := service.GetKey(ctx, tt.keyName)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetKey() error = nil, wantErr %v", tt.wantErr)
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
	repo := newMockKeyRepository()
	logger := zap.NewNop()
	service := NewHLSService(repo, logger)

	ctx := context.Background()
	err := service.ReloadKeys(ctx)
	if err != nil {
		t.Errorf("ReloadKeys() unexpected error = %v", err)
	}
}
