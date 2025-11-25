// Package logger provides structured logging capabilities for the application
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger to provide application-wide logging
type Logger struct {
  *zap.Logger
}

// Config defines logger configuration
type Config struct {
  Level      zapcore.Level
  Encoding   string // json or console
  OutputPath string // stdout, stderr, or file path
}

// New creates a new Logger instance with given configuration
// Returns error if logger initialization fails
func New(cfg Config) (*Logger, error) {
  encoderConfig := zapcore.EncoderConfig{
    TimeKey:        "time",
    LevelKey:       "level",
    NameKey:        "logger",
    CallerKey:      "caller",
    MessageKey:     "msg",
    StacktraceKey:  "stacktrace",
    LineEnding:     zapcore.DefaultLineEnding,
    EncodeLevel:    zapcore.CapitalLevelEncoder,
    EncodeTime:     zapcore.ISO8601TimeEncoder,
    EncodeDuration: zapcore.SecondsDurationEncoder,
    EncodeCaller:   zapcore.ShortCallerEncoder,
  }

  var encoder zapcore.Encoder
  if cfg.Encoding == "json" {
    encoder = zapcore.NewJSONEncoder(encoderConfig)
  } else {
    encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var writeSyncer zapcore.WriteSyncer
	switch cfg.OutputPath {
	case "", "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	default:
		file, err := os.OpenFile(cfg.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}
		writeSyncer = zapcore.AddSync(file)
	}

	core := zapcore.NewCore(encoder, writeSyncer, cfg.Level)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{Logger: zapLogger}, nil
}// NewDevelopment creates a development logger with console output
func NewDevelopment() (*Logger, error) {
  return New(Config{
    Level:      zapcore.DebugLevel,
    Encoding:   "console",
    OutputPath: "stdout",
  })
}

// NewProduction creates a production logger with JSON output
func NewProduction() (*Logger, error) {
  return New(Config{
    Level:      zapcore.InfoLevel,
    Encoding:   "json",
    OutputPath: "stdout",
  })
}
