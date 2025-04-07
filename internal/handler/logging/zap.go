package logging

import (
	//"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	//"time"
	"github.com/gin-gonic/gin"
	"os"
)

var Logger *zap.Logger

func InitZapLogging() *zap.Logger {
	// zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	Logger, _ = zap.NewProduction()
	Logger, _ = zap.NewDevelopment()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:     "time",
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	Logger = zap.New(core)
	defer Logger.Sync()
	return Logger
}

func ZinLogging(c *gin.Context) *zap.Logger {
	// zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	Logger, _ = zap.NewProduction()
	Logger, _ = zap.NewDevelopment()
	Logger, _ = zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
	}.Build()

	return Logger
}

func GinLogging(c *gin.Context) *zap.Logger {
	// zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	Logger, _ = zap.NewProduction()
	Logger, _ = zap.NewDevelopment()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:     "time",
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	Logger = zap.New(core)
	defer Logger.Sync()
	return Logger
}
