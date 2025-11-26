package repository

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"hls-key-server-go/internal/apperrors"
)

func TestNewFileKeyRepository(t *testing.T) {
	tempDir := t.TempDir()
	repo, err := NewFileKeyRepository(tempDir)

	if err != nil {
		t.Fatalf("NewFileKeyRepository() error = %v", err)
	}
	if repo == nil {
		t.Fatal("NewFileKeyRepository() returned nil")
	}
}

func TestFileKeyRepository_Get(t *testing.T) {
	tempDir := t.TempDir()
	keyContent := []byte("test-key-content-1234567890123456")

	// Create a test key file
	keyPath := filepath.Join(tempDir, "test.key")
	if err := os.WriteFile(keyPath, keyContent, 0o644); err != nil {
		t.Fatalf("Failed to create test key: %v", err)
	}

	repo, err := NewFileKeyRepository(tempDir)
	if err != nil {
		t.Fatalf("NewFileKeyRepository() error = %v", err)
	}

	tests := []struct {
		name    string
		keyName string
		wantErr error
		wantLen int
	}{
		{
			name:    "get existing key",
			keyName: "test.key",
			wantErr: nil,
			wantLen: len(keyContent),
		},
		{
			name:    "get non-existent key",
			keyName: "nonexistent.key",
			wantErr: apperrors.ErrKeyNotFound,
			wantLen: 0,
		},
		{
			name:    "empty key name",
			keyName: "",
			wantErr: apperrors.ErrInvalidKeyName,
			wantLen: 0,
		},
		{
			name:    "path traversal attempt with ..",
			keyName: "../etc/passwd.key",
			wantErr: apperrors.ErrInvalidKeyName,
			wantLen: 0,
		},
		{
			name:    "path traversal with forward slash",
			keyName: "dir/test.key",
			wantErr: apperrors.ErrInvalidKeyName,
			wantLen: 0,
		},
		{
			name:    "path traversal with backslash",
			keyName: "dir\\test.key",
			wantErr: apperrors.ErrInvalidKeyName,
			wantLen: 0,
		},
		{
			name:    "missing .key extension",
			keyName: "test.txt",
			wantErr: apperrors.ErrInvalidKeyName,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			got, err := repo.Get(ctx, tt.keyName)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Get() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Get() unexpected error = %v", err)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("Get() returned %d bytes, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestFileKeyRepository_List(t *testing.T) {
	tempDir := t.TempDir()

	// Create multiple test key files
	testKeys := []string{"key1.key", "key2.key", "stream.key"}
	for _, keyName := range testKeys {
		keyPath := filepath.Join(tempDir, keyName)
		if err := os.WriteFile(keyPath, []byte("test-content"), 0o644); err != nil {
			t.Fatalf("Failed to create test key %s: %v", keyName, err)
		}
	}

	repo, err := NewFileKeyRepository(tempDir)
	if err != nil {
		t.Fatalf("NewFileKeyRepository() error = %v", err)
	}

	ctx := context.Background()
	keys := repo.List(ctx)

	if len(keys) != len(testKeys) {
		t.Errorf("List() returned %d keys, want %d", len(keys), len(testKeys))
	}

	// Verify all keys are present
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}
	for _, expectedKey := range testKeys {
		if !keyMap[expectedKey] {
			t.Errorf("List() missing expected key %s", expectedKey)
		}
	}
}

func TestFileKeyRepository_Reload(t *testing.T) {
	tempDir := t.TempDir()

	// Create initial key
	keyPath := filepath.Join(tempDir, "initial.key")
	if err := os.WriteFile(keyPath, []byte("initial"), 0o644); err != nil {
		t.Fatalf("Failed to create initial key: %v", err)
	}

	repo, err := NewFileKeyRepository(tempDir)
	if err != nil {
		t.Fatalf("NewFileKeyRepository() error = %v", err)
	}

	// Add a new key after initialization
	newKeyPath := filepath.Join(tempDir, "new.key")
	if err := os.WriteFile(newKeyPath, []byte("new-content"), 0o644); err != nil {
		t.Fatalf("Failed to create new key: %v", err)
	}

	// Reload should pick up the new key
	ctx := context.Background()
	if err := repo.Reload(ctx); err != nil {
		t.Fatalf("Reload() error = %v", err)
	}

	keys := repo.List(ctx)

	if len(keys) != 2 {
		t.Errorf("After Reload(), List() returned %d keys, want 2", len(keys))
	}

	// Verify the new key is accessible
	data, err := repo.Get(ctx, "new.key")
	if err != nil {
		t.Errorf("Get(new.key) after Reload() error = %v", err)
	}
	if string(data) != "new-content" {
		t.Errorf("Get(new.key) = %s, want new-content", string(data))
	}
}

func BenchmarkFileKeyRepository_Get(b *testing.B) {
	tempDir := b.TempDir()
	keyContent := make([]byte, 16) // Typical AES-128 key size
	for i := range keyContent {
		keyContent[i] = byte(i)
	}

	keyPath := filepath.Join(tempDir, "bench.key")
	if err := os.WriteFile(keyPath, keyContent, 0o644); err != nil {
		b.Fatal(err)
	}

	repo, err := NewFileKeyRepository(tempDir)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Get(ctx, "bench.key")
	}
}
