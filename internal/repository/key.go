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

// validateKeyName performs security validation on key filenames
// to prevent directory traversal attacks and enforce naming conventions.
//
// Requirements:
//   - Must end with .key extension
//   - Cannot be empty
//   - Cannot contain control characters (ASCII 0-31, 127)
//   - Cannot contain path separators (/, \)
//   - Cannot contain parent directory references (..)
//   - Must remain unchanged after filepath.Clean (no manipulation)
func validateKeyName(name string) error {
	// Basic validation
	if name == "" {
		return apperrors.ErrInvalidKeyName
	}
	if !strings.HasSuffix(name, ".key") {
		return apperrors.ErrInvalidKeyName
	}

	// Check for control characters (null byte, newline, etc.)
	for _, ch := range name {
		if ch < 32 || ch == 127 { // ASCII control characters and DEL
			return apperrors.ErrInvalidKeyName
		}
	}

	// Check for malicious characters before cleaning
	if strings.Contains(name, "..") ||
		strings.Contains(name, "/") ||
		strings.Contains(name, "\\") {
		return apperrors.ErrInvalidKeyName
	}

	// Clean and verify path hasn't changed
	cleanName := filepath.Clean(name)
	if cleanName != name {
		return apperrors.ErrInvalidKeyName
	}

	// Ensure no directory separators remain after cleaning
	if strings.Contains(cleanName, string(filepath.Separator)) {
		return apperrors.ErrInvalidKeyName
	}

	return nil
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
	// Validate key name to prevent directory traversal
	if err := validateKeyName(name); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	key, exists := r.cache[name]
	if !exists {
		return nil, apperrors.ErrKeyNotFound
	}

	// Return a copy to prevent cache mutation
	return append([]byte(nil), key...), nil
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
		if file.IsDir() {
			continue
		}

		fileName := file.Name()

		// Apply same validation to loaded files
		if err := validateKeyName(fileName); err != nil {
			// Skip invalid files but continue loading other keys
			continue
		}

		keyPath := filepath.Join(r.keyDir, fileName)
		keyData, err := os.ReadFile(keyPath)
		if err != nil {
			return fmt.Errorf("read key file %s: %w", fileName, err)
		}

		newCache[fileName] = keyData
	}

	r.mu.Lock()
	r.cache = newCache
	r.mu.Unlock()

	return nil
}
