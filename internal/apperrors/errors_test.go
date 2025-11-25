package apperrors

import (
	"errors"
	"testing"
)

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")

	tests := []struct {
		name    string
		err     error
		msg     string
		wantNil bool
	}{
		{
			name:    "wrap non-nil error",
			err:     originalErr,
			msg:     "context",
			wantNil: false,
		},
		{
			name:    "wrap nil error",
			err:     nil,
			msg:     "context",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := Wrap(tt.err, tt.msg)
			if tt.wantNil {
				if wrapped != nil {
					t.Errorf("Wrap() = %v, want nil", wrapped)
				}
				return
			}
			if wrapped == nil {
				t.Error("Wrap() returned nil for non-nil error")
				return
			}
			// Verify the error can be unwrapped
			if !errors.Is(wrapped, originalErr) {
				t.Error("Wrap() error is not unwrappable to original error")
			}
		})
	}
}

func TestWrapf(t *testing.T) {
	originalErr := errors.New("original error")
	wrapped := Wrapf(originalErr, "context %s", "message")

	if wrapped == nil {
		t.Fatal("Wrapf() returned nil")
	}
	if !errors.Is(wrapped, originalErr) {
		t.Error("Wrapf() error is not unwrappable to original error")
	}

	// Test nil error
	wrappedNil := Wrapf(nil, "context")
	if wrappedNil != nil {
		t.Errorf("Wrapf() with nil error = %v, want nil", wrappedNil)
	}
}

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrKeyNotFound", ErrKeyNotFound},
		{"ErrInvalidKeyName", ErrInvalidKeyName},
		{"ErrKeyReadFailed", ErrKeyReadFailed},
		{"ErrUnauthorized", ErrUnauthorized},
		{"ErrTokenInvalid", ErrTokenInvalid},
		{"ErrTokenMissing", ErrTokenMissing},
		{"ErrInvalidCredentials", ErrInvalidCredentials},
		{"ErrMissingHeader", ErrMissingHeader},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s is nil", tt.name)
			}
			// Test that errors can be compared
			if !errors.Is(tt.err, tt.err) {
				t.Errorf("errors.Is(%s, %s) = false", tt.name, tt.name)
			}
			// Test wrapped error can be unwrapped
			wrapped := Wrap(tt.err, "context")
			if !errors.Is(wrapped, tt.err) {
				t.Errorf("wrapped %s cannot be unwrapped", tt.name)
			}
		})
	}
}
