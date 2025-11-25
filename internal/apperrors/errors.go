// Package apperrors defines application-specific error types and helpers
package apperrors

import (
	"errors"
	"fmt"
)

// Domain errors using sentinel pattern
var (
	// ErrKeyNotFound indicates the requested key file does not exist
	ErrKeyNotFound = errors.New("key file not found")

	// ErrInvalidKeyName indicates the key name format is invalid
	ErrInvalidKeyName = errors.New("invalid key file name")

	// ErrKeyReadFailed indicates failure to read key file
	ErrKeyReadFailed = errors.New("failed to read key file")

	// ErrUnauthorized indicates authentication failure
	ErrUnauthorized = errors.New("unauthorized")

	// ErrTokenInvalid indicates JWT token validation failure
	ErrTokenInvalid = errors.New("invalid or expired token")

	// ErrTokenMissing indicates JWT token is not provided
	ErrTokenMissing = errors.New("token is required")

	// ErrInvalidCredentials indicates username or password is incorrect
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrMissingHeader indicates required HTTP header is missing
	ErrMissingHeader = errors.New("required header is missing")
)

// Wrap wraps an error with additional context
// Use this for error propagation across layers
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// Wrapf wraps an error with formatted context
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %w", msg, err)
}

// IsKeyNotFound checks if error is ErrKeyNotFound
func IsKeyNotFound(err error) bool {
	return errors.Is(err, ErrKeyNotFound)
}

// IsInvalidKeyName checks if error is ErrInvalidKeyName
func IsInvalidKeyName(err error) bool {
	return errors.Is(err, ErrInvalidKeyName)
}

// IsInvalidCredentials checks if error is ErrInvalidCredentials
func IsInvalidCredentials(err error) bool {
	return errors.Is(err, ErrInvalidCredentials)
}
