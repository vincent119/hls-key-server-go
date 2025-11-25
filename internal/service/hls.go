// Package service provides business logic layer
package service

import (
	"context"

	"go.uber.org/zap"

	"hls-key-server-go/internal/apperrors"
	"hls-key-server-go/internal/repository"
)

// HLSService handles HLS key business logic
type HLSService struct {
	keyRepo repository.KeyRepository
	logger  *zap.Logger
}

// NewHLSService creates a new HLS service instance
func NewHLSService(keyRepo repository.KeyRepository, logger *zap.Logger) *HLSService {
	return &HLSService{
		keyRepo: keyRepo,
		logger:  logger,
	}
}

// GetKey retrieves an encryption key by name
func (s *HLSService) GetKey(ctx context.Context, keyName string) ([]byte, error) {
	key, err := s.keyRepo.Get(ctx, keyName)
	if err != nil {
		return nil, apperrors.Wrap(err, "get key from repository")
	}

	s.logger.Info("key retrieved",
		zap.String("key_name", keyName),
		zap.Int("key_size", len(key)),
	)

	return key, nil
}

// ListKeys returns all available key names
func (s *HLSService) ListKeys(ctx context.Context) []string {
	return s.keyRepo.List(ctx)
}

// ReloadKeys reloads all keys from storage
func (s *HLSService) ReloadKeys(ctx context.Context) error {
	if err := s.keyRepo.Reload(ctx); err != nil {
		return apperrors.Wrap(err, "reload keys")
	}

	s.logger.Info("keys reloaded successfully")
	return nil
}
