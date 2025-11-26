package repository

import (
	"testing"

	"hls-key-server-go/internal/apperrors"
)

func TestValidateKeyName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		keyName string
		wantErr error
	}{
		// Valid cases
		{
			name:    "valid simple key",
			keyName: "stream.key",
			wantErr: nil,
		},
		{
			name:    "valid with numbers",
			keyName: "stream123.key",
			wantErr: nil,
		},
		{
			name:    "valid with underscore",
			keyName: "stream_backup.key",
			wantErr: nil,
		},
		{
			name:    "valid with dash",
			keyName: "stream-prod.key",
			wantErr: nil,
		},
		{
			name:    "valid with dots",
			keyName: "stream.v2.key",
			wantErr: nil,
		},

		// Invalid cases - empty or wrong extension
		{
			name:    "empty name",
			keyName: "",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "missing extension",
			keyName: "stream",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "wrong extension",
			keyName: "stream.txt",
			wantErr: apperrors.ErrInvalidKeyName,
		},

		// Invalid cases - path traversal attempts
		{
			name:    "parent directory traversal",
			keyName: "../stream.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "double parent directory",
			keyName: "../../etc/passwd.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "hidden parent in middle",
			keyName: "dir/../stream.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "parent at end",
			keyName: "stream.key/..",
			wantErr: apperrors.ErrInvalidKeyName,
		},

		// Invalid cases - directory separators
		{
			name:    "absolute path unix",
			keyName: "/etc/passwd.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "subdirectory unix",
			keyName: "subdir/stream.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "backslash separator",
			keyName: "dir\\stream.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "absolute path windows",
			keyName: "C:\\keys\\stream.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "UNC path",
			keyName: "\\\\server\\share\\stream.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},

		// Invalid cases - control characters
		{
			name:    "null byte",
			keyName: "stream\x00.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "newline in name",
			keyName: "stream\n.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "carriage return",
			keyName: "stream\r.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "tab character",
			keyName: "stream\t.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "vertical tab",
			keyName: "stream\v.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "form feed",
			keyName: "stream\f.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "backspace",
			keyName: "stream\b.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "bell character",
			keyName: "stream\a.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "escape character",
			keyName: "stream\x1b.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},
		{
			name:    "DEL character",
			keyName: "stream\x7f.key",
			wantErr: apperrors.ErrInvalidKeyName,
		},

		// Edge cases
		{
			name:    "only extension",
			keyName: ".key",
			wantErr: nil, // This is technically valid
		},
		{
			name:    "double extension",
			keyName: "stream.key.key",
			wantErr: nil,
		},
		{
			name:    "very long valid name",
			keyName: "stream_with_very_long_name_that_is_still_valid_1234567890.key",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateKeyName(tt.keyName)
			if err != tt.wantErr {
				t.Errorf("validateKeyName(%q) error = %v, wantErr %v", tt.keyName, err, tt.wantErr)
			}
		})
	}
}

func BenchmarkValidateKeyName(b *testing.B) {
	testCases := []string{
		"stream.key",
		"stream123.key",
		"../etc/passwd.key",
		"/absolute/path.key",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validateKeyName(testCases[i%len(testCases)])
	}
}
