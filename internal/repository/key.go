// Package repository provides data access layer for key storage
package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"hls-key-server-go/internal/apperrors"
)

// KeyRepository defines the interface for key storage operations
type KeyRepository interface {
	// Get retrieves a key by name
	Get(ctx context.Context, name string) ([]byte, error)
	// List returns all available key names
	List(ctx context.Context) []string
	// Reload reloads all keys from storage
	Reload(ctx context.Context) error
}

// FileKeyRepository implements KeyRepository using filesystem storage
type FileKeyRepository struct {
	keyDir string
	cache  map[string][]byte
	mu     sync.RWMutex
}

// NewFileKeyRepository creates a new file-based key repository
// keyDir should be an absolute or relative path to the directory containing .key files
func NewFileKeyRepository(keyDir string) (*FileKeyRepository, error) {
	if keyDir == "" {
		return nil, fmt.Errorf("keyDir cannot be empty")
	}

	repo := &FileKeyRepository{
		keyDir: keyDir,
		cache:  make(map[string][]byte),
	}

	// Create directory if not exists
	if err := os.MkdirAll(keyDir, 0o755); err != nil {
		return nil, fmt.Errorf("create key directory: %w", err)
	}

	// Load keys on initialization
	if err := repo.Reload(context.Background()); err != nil {
		return nil, fmt.Errorf("initial key load: %w", err)
	}

	return repo, nil
}

// Get retrieves a key by name from cache
func (r *FileKeyRepository) Get(_ context.Context, name string) ([]byte, error) {
	if name == "" || !strings.HasSuffix(name, ".key") {
		return nil, apperrors.ErrInvalidKeyName
	}

	// Clean path to prevent directory traversal
	name = filepath.Clean(name)
	if strings.Contains(name, "..") {
		return nil, apperrors.ErrInvalidKeyName
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	key, exists := r.cache[name]
	if !exists {
		return nil, apperrors.ErrKeyNotFound
	}

	// Return a copy to prevent cache mutation
	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)

	return keyCopy, nil
}

// List returns all available key names
func (r *FileKeyRepository) List(_ context.Context) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.cache))
	for name := range r.cache {
		names = append(names, name)
	}
	return names
}

// Reload reloads all keys from the filesystem
func (r *FileKeyRepository) Reload(_ context.Context) error {
	files, err := os.ReadDir(r.keyDir)
	if err != nil {
		return fmt.Errorf("read key directory: %w", err)
	}

	newCache := make(map[string][]byte)

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".key") {
			continue
		}

		keyPath := filepath.Join(r.keyDir, file.Name())
		keyData, err := os.ReadFile(keyPath)
		if err != nil {
			return fmt.Errorf("read key file %s: %w", file.Name(), err)
		}

		newCache[file.Name()] = keyData
	}

	r.mu.Lock()
	r.cache = newCache
	r.mu.Unlock()

	return nil
}
